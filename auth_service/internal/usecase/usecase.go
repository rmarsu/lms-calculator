package usecase

import (
	"lms-1/auth_service/internal/models"
	"lms-1/auth_service/internal/repository"
)

type AuthUsecase interface {
	Register(data *models.Register) error
	Login(data *models.Login) (string, error)
}

type authUsecase struct {
	repo repository.AuthRepository
}
