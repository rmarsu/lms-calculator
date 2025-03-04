package queue

import (
	"errors"
	"go/ast"
	"go/token"
	"lms-1/internal/domain"
	"lms-1/pkg/calc"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

type Timings struct {
	TimeAdditionMs       int32
	TimeSubtractionMs    int32
	TimeMultiplicationMs int32
	TimeDivisionMs       int32
}

type Queue struct {
	mu          sync.RWMutex
	Expressions map[string]*domain.Expression
	Tasks       map[string]*domain.Task
	Timings
}

func New(t Timings) *Queue {
	return &Queue{
		mu:          sync.RWMutex{},
		Expressions: make(map[string]*domain.Expression),
		Tasks:       make(map[string]*domain.Task),
		Timings:     t,
	}
}

func (q *Queue) AddExpression(expr domain.Expression) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Expressions[expr.Id] = &expr
}

func (q *Queue) GetExpressionById(id string) (*domain.Expression, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	expr, ok := q.Expressions[id]
	return expr, ok
}

func (q *Queue) GetExpressions() []domain.Expression {
	q.mu.RLock()
	defer q.mu.RUnlock()
	var exprs []domain.Expression
	for _, expr := range q.Expressions {
		exprs = append(exprs, *expr)
	}
	return exprs
}

func (q *Queue) GetExpression() (*domain.Expression, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if len(q.Expressions) == 0 {
		return nil, false
	}
	for _, expr := range q.Expressions {
		return expr, true
	}
	return nil, false
}

func (q *Queue) GetTask() (*domain.Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.Tasks) == 0 {
		return nil, false
	}
	for _, task := range q.Tasks {
		return task, true
	}
	return nil, false
}

func (q *Queue) RunTask(id string) error {
	q.mu.Lock()
	expression, ok := q.GetExpressionById(id)
	q.mu.Unlock()
	if !ok {
		return errors.New(domain.ErrInvalidExpression)
	}
	expression.Status = domain.StatusInProgress

	ast, err := calc.ParseAst(expression.Expression)
	if err != nil {
		return errors.New(domain.ErrInvalidExpression)
	}

	resChan := make(chan float64)

	go func() {
		defer close(resChan)
		result := q.evaluateAst(ast, resChan)
		expression.Status = domain.StatusCompleted
		expression.Result = result
	}()

	return nil
}

func (q *Queue) RemoveTask(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.Tasks, id)
}

func (q *Queue) RollbackResult(id string, res float64) error {
	q.mu.RLock()
	task, ok := q.Tasks[id]
	q.mu.RUnlock()
	if !ok {
		return errors.New(domain.ErrTaskNotFound)
	}

	task.ResultChan <- res
	return nil
}

func (q *Queue) evaluateAst(node ast.Node, res chan float64) float64 {
	switch n := node.(type) {
	case *ast.ParenExpr:
		return q.evaluateAst(n.X, res)
	case *ast.BinaryExpr:
		left := q.evaluateAst(n.X, res)
		right := q.evaluateAst(n.Y, res)

		var timing int32
		switch n.Op {
		case token.ADD:
			timing = q.TimeAdditionMs
		case token.SUB:
			timing = q.TimeSubtractionMs
		case token.MUL:
			timing = q.TimeMultiplicationMs
		case token.QUO:
			timing = q.TimeDivisionMs
		}

		task := &domain.Task{
			Id:              uuid.NewString(),
			Arg1:            left,
			Arg2:            right,
			Operation:       n.Op.String(),
			OperationTimeMs: timing,
			ResultChan:      res,
		}
		q.Tasks[task.Id] = task

		result := <-res
		q.RemoveTask(task.Id)
		return result
	case *ast.BasicLit:
		val, _ := strconv.ParseFloat(n.Value, 64)
		return val
	default:
		return 0.0
	}
}

func (q *Queue) RemoveExpression(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.Expressions, id)
}

func (q *Queue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.Expressions)
}
