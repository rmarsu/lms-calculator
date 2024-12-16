package app

import (
	"context"
	"lms-1/internal/config"
	"lms-1/internal/server"
	"lms-1/internal/service"
	transport "lms-1/internal/transport/http"
	"lms-1/pkg/calc"
	"lms-1/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	configPath = "configs/config.yaml"
)

func Run() {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	services := service.NewService(&service.Deps{
		Calc: calc.NewCalc(),
	})

	handler := transport.NewHandler(services)

	srv := server.NewServer(cfg, handler.InitRoutes())

	go func() {
		if err := srv.Run(); err != nil {
			logger.Error(err)
		}
	}()

	logger.Infof("Сервер слушает на... :%s", cfg.HTTP.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("Попытка остановить сервер закончилась неудачно: %v", err)
	}
}
