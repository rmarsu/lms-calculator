package agent

import (
	"context"
	"errors"
	"lms-1/agent_service/internal/models"
	pb_orchestrator "lms-1/pkg/orchestrator"
	"time"

	"go.uber.org/zap"
)

var (
	ErrDivisionByZero = errors.New("division by zero")
	ErrUnknownOp      = errors.New("unknown op")
)

func evaluateExpression(task *models.Task) (float32, error) {
	switch task.Op {
	case "+":
		return task.Arg1 + task.Arg2, nil
	case "-":
		return task.Arg1 - task.Arg2, nil
	case "*":
		return task.Arg1 * task.Arg2, nil
	case "/":
		if task.Arg2 == 0 {
			return 0.0, ErrDivisionByZero
		}
		return task.Arg1 / task.Arg2, nil
	}
	return 0.0, ErrUnknownOp
}

func RunAgent(client pb_orchestrator.OrchestratorServiceClient, sugar *zap.SugaredLogger) {
	for {
		task, err := client.GiveTask(context.Background(), &pb_orchestrator.Empty{})
		if err != nil {
			sugar.Errorf("error while receiving task: %v", err)
			time.Sleep(time.Second)
			continue
		}

		result, err := evaluateExpression(&models.Task{
			Id:       task.GetId(),
			Arg1:     task.GetArg1(),
			Arg2:     task.GetArg2(),
			Op:       task.GetOp(),
			OpTimeMs: task.GetOpTimeMs(),
		})
		if err != nil {
			sugar.Errorf("error during evaluation: %v", err)
			continue
		}

		_, err = client.AnswerTask(context.Background(), &pb_orchestrator.Answer{
			Id:  task.GetId(),
			Res: result,
		})
		if err != nil {
			sugar.Errorf("error while sending the answer: %v", err)
		} else {
			sugar.Infof("answer sent: %f", result)
		}

		time.Sleep(time.Duration(task.OpTimeMs) * time.Millisecond)
	}
}
