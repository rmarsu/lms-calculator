package repository

import (
	"database/sql"
	"lms-1/orchestrator_service/internal/models"
)

type sqliteExpressionsRepository struct {
	conn *sql.DB
}

func NewSqliteExpressionsRepository(conn *sql.DB) ExpressionsRepository {
	return &sqliteExpressionsRepository{
		conn: conn,
	}
}

func (r *sqliteExpressionsRepository) Add(expr *models.Expression) (int64, error) {
	query := `
		INSERT INTO expressions (expression, status, result)
		VALUES (?, ?, ?)
	`
	result, err := r.conn.Exec(query, expr.Expression, expr.Status, expr.Result)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *sqliteExpressionsRepository) Update(id int64, status models.Status, result float64) error {
	query := `
		UPDATE expressions
		SET status = ?, result = ?
		WHERE id = ?
	`
	_, err := r.conn.Exec(query, status, result, id)
	return err
}

func (r *sqliteExpressionsRepository) GetExprById(id int64) (*models.Expression, error) {
	query := `
		SELECT id, expression, status, result
		FROM expressions
		WHERE id = ?
	`
	row := r.conn.QueryRow(query, id)

	var expr models.Expression
	err := row.Scan(&expr.Id, &expr.Expression, &expr.Status, &expr.Result)
	if err != nil {
		return nil, err
	}
	return &expr, nil
}

func (r *sqliteExpressionsRepository) GetExpressions() ([]models.Expression, error) {
	query := `
		SELECT id, expression, status, result
		FROM expressions
	`
	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []models.Expression
	for rows.Next() {
		var expr models.Expression
		if err := rows.Scan(&expr.Id, &expr.Expression, &expr.Status, &expr.Result); err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expressions, nil
}
