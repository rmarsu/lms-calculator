package server

import (
	"lms-1/auth_service/internal/usecase"
)

type server struct {
	usecase usecase.AuthUsecase
}

func New(uc usecase.AuthUsecase) *server {
	return &server{
		usecase: uc,
	}
}



