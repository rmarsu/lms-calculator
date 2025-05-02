package server

import (
	"context"
	"lms-1/orchestrator_service/internal/models"
	"lms-1/orchestrator_service/internal/usecase"
	pb_orchestrator "lms-1/pkg/orchestrator"
)

var _ pb_orchestrator.OrchestratorServiceServer = (*server)(nil)

type server struct {
	pb_orchestrator.UnimplementedOrchestratorServiceServer
	usecase usecase.OrchestratorUsecase
}

func New(u usecase.OrchestratorUsecase) pb_orchestrator.OrchestratorServiceServer {
	return &server{
		usecase: u,
	}
}

func (s *server) GetExpressions(_ *pb_orchestrator.Empty, stream pb_orchestrator.OrchestratorService_GetExpressionsServer) error {
	expressions, err := s.usecase.GetExpressions()
	if err != nil {
		return err
	}

	for _, expr := range expressions {
		if err := stream.Send(&pb_orchestrator.Expression{
			Id:     expr.Id,
			Status: statusToPbStatus(expr.Status),
			Result: float32(expr.Result),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) GetExpressionById(ctx context.Context, req *pb_orchestrator.Id) (*pb_orchestrator.Expression, error) {
	expr, err := s.usecase.GetExpressionById(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb_orchestrator.Expression{
		Id:     expr.Id,
		Status: statusToPbStatus(expr.Status),
		Result: float32(expr.Result),
	}, nil
}

func (s *server) AddToQueue(ctx context.Context, req *pb_orchestrator.StringExpression) (*pb_orchestrator.Expression, error) {
	expr := &models.Expression{
		Expression: req.Expr,
		Status:     models.StatusPending,
	}

	id, err := s.usecase.RunExpression(expr)
	if err != nil {
		return nil, err
	}

	return &pb_orchestrator.Expression{
		Id:     id,
		Status: statusToPbStatus(expr.Status),
		Result: float32(expr.Result),
	}, nil
}

func (s *server) GiveTask(ctx context.Context, _ *pb_orchestrator.Empty) (*pb_orchestrator.Task, error) {
	task, err := s.usecase.GiveTask()
	if err != nil {
		return nil, err
	}

	return &pb_orchestrator.Task{
		Id:       task.Id,
		Arg1:     float32(task.Arg1),
		Arg2:     float32(task.Arg2),
		Op:       task.Operation,
		OpTimeMs: task.OperationTime,
	}, nil
}

func (s *server) AnswerTask(ctx context.Context, req *pb_orchestrator.Answer) (*pb_orchestrator.Empty, error) {
	err := s.usecase.AnswerTask(req.Id, float64(req.Res))
	if err != nil {
		return nil, err
	}

	return &pb_orchestrator.Empty{}, nil
}

func statusToPbStatus(status models.Status) pb_orchestrator.Status {
	switch status {
	case models.StatusInProgress:
		return pb_orchestrator.Status_StatusInProgress
	case models.StatusCompleted:
		return pb_orchestrator.Status_StatusCompleted
	case models.StatusPending:
		return pb_orchestrator.Status_StatusPending
	}
	return pb_orchestrator.Status_StatusInProgress
}
