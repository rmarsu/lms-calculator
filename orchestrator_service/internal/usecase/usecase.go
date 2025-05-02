package usecase

import (
	"go/ast"
	"go/token"
	"lms-1/orchestrator_service/internal/models"
	"lms-1/orchestrator_service/internal/repository"
	"lms-1/pkg/calc"
	"strconv"
	"sync"

	"go.uber.org/zap"
)

type OrchestratorUsecase interface {
	RunExpression(expr *models.Expression) (int64, error)
	GetExpressions() ([]models.Expression, error)
	GetExpressionById(id int64) (*models.Expression, error)
	GiveTask() (*models.Task, error)
	AnswerTask(id int64, res float64) error
}

type Timings struct {
	TimeAdditionMs       int64
	TimeSubtractionMs    int64
	TimeMultiplicationMs int64
	TimeDivisionMs       int64
}

type orchestratorUsecase struct {
	mu         sync.Mutex
	expr_repo  repository.ExpressionsRepository
	tasks_repo repository.TasksRepository
	Timings
	sugar    *zap.SugaredLogger
	chansMap map[int64]chan float64
}

func NewOrchestratorUsecase(expr_repo repository.ExpressionsRepository, tasks_repo repository.TasksRepository, t Timings, sugar *zap.SugaredLogger) OrchestratorUsecase {
	return &orchestratorUsecase{
		mu:         sync.Mutex{},
		expr_repo:  expr_repo,
		tasks_repo: tasks_repo,
		Timings:    t,
		sugar:      sugar,
		chansMap:   make(map[int64]chan float64),
	}
}

func (uc *orchestratorUsecase) AddTask(task *models.Task) (chan float64, error) {
	uc.sugar.Infow("Adding task", "task", task)

	id, err := uc.tasks_repo.Add(task)
	if err != nil {
		uc.sugar.Errorw("Failed to add task to repository", "error", err)
		return nil, ErrDeadDB
	}

	task.Id = id
	resChan := make(chan float64)

	uc.mu.Lock()
	uc.chansMap[id] = resChan
	uc.mu.Unlock()

	uc.sugar.Infow("Task added successfully", "task_id", id)
	return resChan, nil
}

func (uc *orchestratorUsecase) DeleteTask(id int64) error {
	uc.sugar.Infow("Deleting task", "task_id", id)

	err := uc.tasks_repo.Delete(id)
	if err != nil {
		uc.sugar.Errorw("Failed to delete task from repository", "error", err)
		return ErrDeadDB
	}

	uc.mu.Lock()
	delete(uc.chansMap, id)
	uc.mu.Unlock()

	uc.sugar.Infow("Task deleted successfully", "task_id", id)
	return nil
}

func (uc *orchestratorUsecase) RunExpression(expr *models.Expression) (int64, error) {
	uc.sugar.Infow("Running expression", "expression", expr.Expression)

	uc.mu.Lock()
	defer uc.mu.Unlock()

	id, err := uc.expr_repo.Add(expr)
	if err != nil {
		uc.sugar.Errorw("Failed to add expression to repository", "error", err)
		return 0, ErrDeadDB
	}

	err = uc.expr_repo.Update(id, models.StatusInProgress, expr.Result)
	if err != nil {
		uc.sugar.Errorw("Failed to update expression status", "error", err)
		return 0, ErrDeadDB
	}

	ast, err := calc.ParseAst(expr.Expression)
	if err != nil {
		uc.sugar.Errorw("Failed to parse expression AST", "error", err)
		return 0, ErrInvalidExpression
	}

	go func() {
		res, err := uc.evaluateAst(ast)
		if err != nil {
			uc.sugar.Errorw("Error during AST evaluation", "error", err)
			return
		}
		err = uc.expr_repo.Update(id, models.StatusCompleted, res)
		if err != nil {
			uc.sugar.Errorw("Failed to update expression result", "error", err)
			return
		}
		uc.sugar.Infow("Expression evaluated successfully", "expression_id", id, "result", res)
	}()

	return id, nil
}

func (uc *orchestratorUsecase) AnswerTask(id int64, res float64) error {
	uc.mu.Lock()
	ch, ok := uc.chansMap[id]
	uc.mu.Unlock()

	if !ok {
		uc.sugar.Warnw("Attempted to answer non-existent or already-answered task", "task_id", id)
		return ErrTaskDoNotExist
	}

	defer func() {
		uc.mu.Lock()
		close(ch)
		delete(uc.chansMap, id)
		uc.mu.Unlock()
	}()

	ch <- res
	return nil
}


func (uc *orchestratorUsecase) GiveTask() (*models.Task, error) {
	uc.sugar.Info("Fetching next task")
	tsk, err := uc.tasks_repo.GetByDate()
	if err != nil {
		uc.sugar.Errorw("Failed to fetch task", "error", err)
		return nil, ErrDeadDB
	}
	uc.sugar.Infow("Task fetched", "task", tsk)
	return tsk, nil
}

func (uc *orchestratorUsecase) GetExpressions() ([]models.Expression, error) {
	uc.sugar.Info("Fetching all expressions")
	exprs, err := uc.expr_repo.GetExpressions()
	if err != nil {
		uc.sugar.Errorw("Failed to fetch expressions", "error", err)
		return nil, ErrDeadDB
	}
	return exprs, nil
}

func (uc *orchestratorUsecase) GetExpressionById(id int64) (*models.Expression, error) {
	uc.sugar.Infow("Fetching expression by ID", "expression_id", id)

	expr, err := uc.expr_repo.GetExprById(id)
	if err != nil {
		uc.sugar.Errorw("Failed to fetch expression", "error", err, "expression_id", id)
		return nil, ErrDeadDB
	}
	return expr, nil
}

func (uc *orchestratorUsecase) evaluateAst(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.ParenExpr:
		return uc.evaluateAst(n.X)
	case *ast.BinaryExpr:
		left, err := uc.evaluateAst(n.X)
		if err != nil {
			return 0.0, err
		}
		right, err := uc.evaluateAst(n.Y)
		if err != nil {
			return 0.0, err
		}

		var timing int64
		switch n.Op {
		case token.ADD:
			timing = uc.TimeAdditionMs
		case token.SUB:
			timing = uc.TimeSubtractionMs
		case token.MUL:
			timing = uc.TimeMultiplicationMs
		case token.QUO:
			timing = uc.TimeDivisionMs
		}

		task := &models.Task{
			Id:            0,
			Arg1:          left,
			Arg2:          right,
			Operation:     n.Op.String(),
			OperationTime: timing,
		}
		res, err := uc.AddTask(task)
		if err != nil {
			return 0.0, err
		}
		result := <-res
		uc.DeleteTask(task.Id)
		return result, nil
	case *ast.BasicLit:
		val, _ := strconv.ParseFloat(n.Value, 64)
		return val, nil
	default:
		return 0.0, ErrInvalidExpression
	}
}
