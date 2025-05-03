package usecase

import "errors"

var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrDeadDB            = errors.New("dead db")
	ErrTaskDoNotExist    = errors.New("task do not exists")
	ErrExprNotFound      = errors.New("expression not found")
)
