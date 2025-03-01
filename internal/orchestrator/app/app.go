package orchestrator_app

import (
	"context"
	"lms-1/internal/orchestrator/server"
	"lms-1/internal/orchestrator/service"
	"lms-1/internal/orchestrator/transport/http"
	"lms-1/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func Run() {
	logger.Info("orchestrator started")

	mustInitEnv()

	service := service.New(&service.Deps{
		TimeAdditionMs:       1000,
		TimeSubtractionMs:    500,
		TimeMultiplicationMs: 2000,
		TimeDivisionMs:       1500,
	})

	handlers := http.NewHandler(service)

	server := server.NewServer(server.NewDefaultServerConfig(), handlers.InitRoutes())
	go func() {
		if err := server.Run(); err != nil {
			logger.Error(err)
		}
	}()

	logger.Infof("orchestrator listening on port %s", server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logger.Error(err)
	}
}

func mustInitEnv() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	logger.Info("environment variables loaded")
}
