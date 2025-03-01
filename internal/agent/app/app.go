package agent_app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lms-1/internal/domain"
	"lms-1/pkg/logger"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func Run() {
	logger.Info("agent started")
	mustInitEnv()

	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	orchestratorPort := os.Getenv("ORCHESTRATOR_PORT")
	if err != nil || computingPower <= 0 {
		logger.Error("invalid computing power")
		os.Exit(1)
	}
	for i := 0; i < computingPower; i++ {
		go worker(orchestratorPort)
		time.Sleep(time.Millisecond * 10)
		logger.Infof("created worker %d", i+1)
	}
}

func worker(port string) {
	for {
		task, err := fetchTask(port)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		result, err := solveTask(task)
		if err != nil {
			logger.Errorf("failed to solve task: %v", err)
			continue
		}

		err = sendResult(port, task.Id, result)
		if err != nil {
			logger.Errorf("failed to send result: %v", err)
			continue
		}

		logger.Infof("solved task %s: result=%f", task.Id, result)
	}
}

func fetchTask(port string) (*domain.Task, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/internal/task", port))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch task: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var task domain.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func solveTask(task *domain.Task) (float64, error) {
	logger.Infof("Solving task %2.f%s%2.f", task.Arg1, task.Operation, task.Arg2)
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2, nil
	case "-":
		return task.Arg1 - task.Arg2, nil
	case "*":
		return task.Arg1 * task.Arg2, nil
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return task.Arg1 / task.Arg2, nil
	default:
		return 0, fmt.Errorf("unsupported operation: %s", task.Operation)
	}
}

func sendResult(port, taskId string, result float64) error {
	data := map[string]interface{}{
		"id":     taskId,
		"result": result,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/internal/task", port), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send result: %s", string(body))
	}

	return nil
}

func mustInitEnv() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	logger.Info("environment variables loaded")
}
