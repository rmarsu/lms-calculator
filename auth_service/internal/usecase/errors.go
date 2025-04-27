package usecase

import "errors"

var (
	ErrPasswordIsTooShort = errors.New("password is too short")
	ErrAlreadyExists      = errors.New("already exists")
	ErrUnknown            = errors.New("unknown error")
	ErrNotFound           = errors.New("not found")
	ErrFailedHash         = errors.New("failed to hash password")
	ErrUnauthenticated    = errors.New("unauthenticated")
)
