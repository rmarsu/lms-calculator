package server

import (
	"context"
	"lms-1/auth_service/internal/models"
	"lms-1/auth_service/internal/usecase"
	pb_auth "lms-1/pkg/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb_auth.AuthServiceServer
	usecase usecase.AuthUsecase
}

func New(uc usecase.AuthUsecase) *server {
	return &server{
		usecase: uc,
	}
}

func (s *server) Register(ctx context.Context,
	identity *pb_auth.Identity) (*pb_auth.Empty, error) {

	if identity.Username == "" || identity.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "no username or password provided")
	}

	err := s.usecase.Register(&models.Register{
		Username: identity.Username,
		Password: identity.Password,
	})
	if err != nil {
		switch err {
		case usecase.ErrAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		case usecase.ErrPasswordIsTooShort:
			return nil, status.Error(codes.Aborted, "password is too short")
		default:
			return nil, status.Error(codes.Internal, "something went wrong")
		}
	}
	return &pb_auth.Empty{}, nil
}

func (s *server) Login(ctx context.Context,
	identity *pb_auth.Identity) (*pb_auth.Token, error) {

	if identity.Username == "" || identity.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "no username or password provided")
	}

	token, err := s.usecase.Login(&models.Login{
		Username: identity.Username,
		Password: identity.Password,
	})
	if err != nil {
		switch err {
		case usecase.ErrNotFound:
			return nil, status.Error(codes.NotFound, "user with this username not found")
		case usecase.ErrUnauthenticated:
			return nil, status.Error(codes.Unauthenticated, "password is incorrect")
		default:
			return nil, status.Error(codes.Internal, "something went wrong")
		}
	}
	return &pb_auth.Token{
		Token: token,
	}, nil
}
