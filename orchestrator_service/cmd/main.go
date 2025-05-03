package main

import (
	_ "github.com/mattn/go-sqlite3"

	"context"
	"fmt"
	"lms-1/orchestrator_service/internal/config"
	"lms-1/orchestrator_service/internal/repository"
	"lms-1/orchestrator_service/internal/server"
	interceptor "lms-1/orchestrator_service/internal/server/interceptors"
	"lms-1/orchestrator_service/internal/usecase"
	"lms-1/pkg/jwt"
	pb_orchestrator "lms-1/pkg/orchestrator"
	"lms-1/pkg/sqlite"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zap logger: %v", err))
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	cfg := config.MustLoad()
	jwtMgr, err := jwt.NewManager(cfg.JwtSecret)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.NewUnaryAuthInterceptor(jwtMgr)),
	)

	conn := sqlite.MustConnect(cfg.SqlitePath)
	defer conn.Close()

	expr_repository := repository.NewSqliteExpressionsRepository(conn)
	tasks_repository := repository.NewSqliteTasksRepository(conn)

	uc := usecase.NewOrchestratorUsecase(expr_repository, tasks_repository, usecase.Timings{
		TimeAdditionMs:       cfg.TimeAdditionMs,
		TimeSubtractionMs:    cfg.TimeSubtractionMs,
		TimeMultiplicationMs: cfg.TimeMultiplicationMs,
		TimeDivisionMs:       cfg.TimeDivisionMs,
	}, sugar)
	pb_orchestrator.RegisterOrchestratorServiceServer(grpcServer, server.New(uc))

	go func() {
		listener, err := net.Listen("tcp", cfg.GrpcPort)
		if err != nil {
			sugar.Fatalw("Failed to listen", "error", err)
		}
		sugar.Infow("gRPC server listening", "port", cfg.GrpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			sugar.Fatalw("Failed to serve gRPC", "error", err)
		}
	}()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb_orchestrator.RegisterOrchestratorServiceHandlerFromEndpoint(ctx, mux, cfg.GrpcPort, opts); err != nil {
		sugar.Fatalw("Failed to register gRPC gateway", "error", err)
	}

	httpServer := &http.Server{
		Addr:    cfg.RestPort,
		Handler: mux,
	}

	go func() {
		sugar.Infow("HTTP server listening", "port", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalw("Failed to serve HTTP", "error", err)
		}
	}()

	<-sigCh
	sugar.Infow("Shutdown signal received")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := httpServer.Shutdown(ctxShutdown); err != nil {
		sugar.Fatalw("HTTP shutdown error", "error", err)
	}

	grpcServer.GracefulStop()
	sugar.Infow("Servers gracefully stopped")
}
