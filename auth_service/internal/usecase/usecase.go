package usecase

import (
	"errors"
	"lms-1/auth_service/internal/models"
	"lms-1/auth_service/internal/repository"
	"lms-1/pkg/hash"
	"lms-1/pkg/jwt"
	"time"
	"unicode/utf8"
)

type AuthUsecase interface {
	Register(data *models.Register) error
	Login(data *models.Login) (string, error)
}

type authUsecase struct {
	repo   repository.AuthRepository
	hasher hash.PasswordHasher
	jwtMgr jwt.TokenManager
	tll    time.Duration
}

func NewAuthUsecase(r repository.AuthRepository, h hash.PasswordHasher, j jwt.TokenManager) AuthUsecase {
	return &authUsecase{
		repo:   r,
		hasher: h,
		jwtMgr: j,
	}
}

func (u *authUsecase) Register(data *models.Register) error {
	if err := isPasswordHard(data.Password); err != nil {
		return err
	}

	hashed_pswrd, err := u.hasher.Hash(data.Password)
	if err != nil {
		return ErrFailedHash
	}

	err = u.repo.Create(&models.User{
		Username:  data.Username,
		Password:  hashed_pswrd,
		CreatedAt: time.Now(),
	})
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrAlreadyExists
		} else {
			return ErrUnknown
		}
	}
	return nil
}

func (u *authUsecase) Login(data *models.Login) (string, error) {
	found_user, err := u.repo.GetByUsername(data.Username)
	if err != nil {
		switch err {
		case repository.ErrNoRows:
			return "", ErrNotFound
		default:
			return "", ErrUnknown
		}
	}
	hashed_pswrd, err := u.hasher.Hash(data.Password)
	if err != nil {
		return "", ErrFailedHash
	}

	if hashed_pswrd != found_user.Password {
		return "", ErrUnauthenticated
	}

	token, err := u.jwtMgr.NewJWT(u.tll)
	if err != nil {
		return "", ErrUnauthenticated
	}
	return token, nil
}

func isPasswordHard(pswrd string) error {
	if utf8.RuneCountInString(pswrd) < 8 {
		return ErrPasswordIsTooShort
	}
	return nil
}
