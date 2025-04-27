package repository

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNoRows        = errors.New("no rows in result set")
)
