package main

import (
	"fmt"
	"log"

	"lms-1/agent_service/internal/agent"
	"lms-1/agent_service/internal/config"
	pb_orchestrator "lms-1/pkg/orchestrator"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoad()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	conn, err := grpc.NewClient(fmt.Sprintf("%s%s", cfg.Host, cfg.OrchestratorPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		sugar.Fatalf("failed to connect to orchestrator: %v", err)
	}
	defer conn.Close()

	client := pb_orchestrator.NewOrchestratorServiceClient(conn)

	numAgents := cfg.ComputingPower

	sugar.Infof("starting %d agents...", numAgents)

	for i := 0; i < numAgents; i++ {
		go agent.RunAgent(client, sugar.With("agent_id", i))
	}

	select {}
}
