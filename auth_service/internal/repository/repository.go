package repository

import "lms-1/auth_service/internal/models"

type AuthRepository interface {
	Create(user *models.User) error
	GetById(id int64) (*models.User, error)
	Delete(id int64) error
}
