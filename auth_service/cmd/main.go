package main

import (
	"context"
	"fmt"
	"lms-1/auth_service/internal/config"
	"lms-1/auth_service/internal/repository"
	"lms-1/auth_service/internal/server"
	"lms-1/auth_service/internal/usecase"
	pb_auth "lms-1/pkg/auth"
	"lms-1/pkg/hash"
	"lms-1/pkg/jwt"
	"lms-1/pkg/sqlite"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/mattn/go-sqlite3"
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

	grpcServer := grpc.NewServer()

	conn := sqlite.MustConnect(cfg.SqlitePath)
	defer conn.Close()

	repository := repository.NewSqliteRepository(conn)
	hasher := hash.NewSHA256Hasher(cfg.HasherSalt)
	jwtMgr, err := jwt.NewManager(cfg.JwtSecret)
	if err != nil {
		sugar.Fatalw("Failed to create JWT manager", "error", err)
	}

	uc := usecase.NewAuthUsecase(repository, hasher, jwtMgr, sugar)
	pb_auth.RegisterAuthServiceServer(grpcServer, server.New(uc))

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
	if err := pb_auth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.GrpcPort, opts); err != nil {
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
