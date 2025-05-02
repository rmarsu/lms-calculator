package repository

import (
	"database/sql"
	"lms-1/orchestrator_service/internal/models"
)

type sqliteTasksRepository struct {
	conn *sql.DB
}

func NewSqliteTasksRepository(conn *sql.DB) TasksRepository {
	return &sqliteTasksRepository{
		conn: conn,
	}
}

func (r *sqliteTasksRepository) Add(task *models.Task) (int64, error) {
	query := `
		INSERT INTO tasks (arg1, arg2, operation, operation_time_ms, expr_id)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.conn.Exec(query, task.Arg1, task.Arg2, task.Operation, task.OperationTime, task.ExprID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *sqliteTasksRepository) GetByDate() (*models.Task, error) {
	query := `
		SELECT id, arg1, arg2, operation, operation_time_ms, expr_id, created_at
		FROM tasks
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := r.conn.QueryRow(query)

	var task models.Task
	err := row.Scan(&task.Id, &task.Arg1, &task.Arg2, &task.Operation, &task.OperationTime, &task.ExprID, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *sqliteTasksRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.conn.Exec(query, id)
	return err
}
