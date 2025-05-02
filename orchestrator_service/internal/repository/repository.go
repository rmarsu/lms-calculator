package repository

import "lms-1/orchestrator_service/internal/models"

type TasksRepository interface {
	Add(task *models.Task) (int64, error)
	GetByDate() (*models.Task, error)
	Delete(id int64) error
}

type ExpressionsRepository interface {
	Add(expr *models.Expression) (int64, error)
	Update(id int64, status models.Status, result float64) error

	GetExprById(id int64) (*models.Expression, error)
	GetExpressions() ([]models.Expression, error)
}
