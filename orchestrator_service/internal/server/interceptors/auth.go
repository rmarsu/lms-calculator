package interceptor

import (
	"context"
	"errors"
	"lms-1/pkg/jwt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewUnaryAuthInterceptor(jwtMgr *jwt.Manager) grpc.UnaryServerInterceptor {
	protectedMethods := map[string]bool{
		"/orchestrator.OrchestratorService/GetExpressions": true,
		"/orchestrator.OrchestratorService/GetExpressionById": true,
		"/orchestrator.OrchestratorService/AddToQueue": true,
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if protectedMethods[info.FullMethod] {
			if err := authorize(ctx, jwtMgr); err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}

func authorize(ctx context.Context, jwtMgr *jwt.Manager) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("missing metadata")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization token not supplied")
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return errors.New("invalid authorization header format")
	}

	tokenStr := parts[1]

	_, err := jwtMgr.Parse(tokenStr)
	if err != nil {
		return errors.New("invalid token")
	}

	return nil
}
