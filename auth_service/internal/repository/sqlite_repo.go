package repository

import (
	"database/sql"
	"errors"
	"lms-1/auth_service/internal/models"
)

type sqliteRepository struct {
	conn *sql.DB
}

func NewSqliteRepository(conn *sql.DB) AuthRepository {
	return &sqliteRepository{
		conn: conn,
	}
}

func (r *sqliteRepository) Create(u *models.User) error {
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	result, err := r.conn.Exec(query, u.Username, u.Password)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = id
	return nil
}

func (r *sqliteRepository) GetById(id int64) (*models.User, error) {
	query := `SELECT id, username, password, timestamp FROM users WHERE id = ?`
	row := r.conn.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *sqliteRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.conn.Exec(query, id)
	return err
}
