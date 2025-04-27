package usecase

import (
	"errors"
	"lms-1/auth_service/internal/models"
	"lms-1/auth_service/internal/repository"
	"lms-1/pkg/hash"
	"lms-1/pkg/jwt"
	"time"
	"unicode/utf8"

	"go.uber.org/zap"
)

type AuthUsecase interface {
	Register(data *models.Register) error
	Login(data *models.Login) (string, error)
}

type authUsecase struct {
	repo   repository.AuthRepository
	hasher hash.PasswordHasher
	jwtMgr jwt.TokenManager
	sugar  *zap.SugaredLogger
	tll    time.Duration
}

func NewAuthUsecase(r repository.AuthRepository, h hash.PasswordHasher, j jwt.TokenManager, s *zap.SugaredLogger) AuthUsecase {
	return &authUsecase{
		repo:   r,
		hasher: h,
		jwtMgr: j,
		sugar:  s,
		tll:    15 * time.Minute,
	}
}

func (u *authUsecase) Register(data *models.Register) error {
	u.sugar.Infow("register attempt", "username", data.Username)

	if err := isPasswordHard(data.Password); err != nil {
		u.sugar.Warnw("password too weak", "username", data.Username, "error", err)
		return err
	}

	hashed_pswrd, err := u.hasher.Hash(data.Password)
	if err != nil {
		u.sugar.Errorw("failed to hash password", "username", data.Username, "error", err)
		return ErrFailedHash
	}

	err = u.repo.Create(&models.User{
		Username:  data.Username,
		Password:  hashed_pswrd,
		CreatedAt: time.Now(),
	})
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			u.sugar.Warnw("user already exists", "username", data.Username)
			return ErrAlreadyExists
		}
		u.sugar.Errorw("failed to create user", "username", data.Username, "error", err)
		return ErrUnknown
	}

	u.sugar.Infow("user registered successfully", "username", data.Username)
	return nil
}

func (u *authUsecase) Login(data *models.Login) (string, error) {
	u.sugar.Infow("login attempt", "username", data.Username)

	found_user, err := u.repo.GetByUsername(data.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			u.sugar.Warnw("user not found", "username", data.Username)
			return "", ErrNotFound
		}
		u.sugar.Errorw("error fetching user", "username", data.Username, "error", err)
		return "", ErrUnknown
	}

	hashed_pswrd, err := u.hasher.Hash(data.Password)
	if err != nil {
		u.sugar.Errorw("failed to hash password on login", "username", data.Username, "error", err)
		return "", ErrFailedHash
	}

	if hashed_pswrd != found_user.Password {
		u.sugar.Warnw("invalid credentials", "username", data.Username)
		return "", ErrUnauthenticated
	}

	token, err := u.jwtMgr.NewJWT(u.tll)
	if err != nil {
		u.sugar.Errorw("failed to generate JWT", "username", data.Username, "error", err)
		return "", ErrUnauthenticated
	}

	u.sugar.Infow("login successful", "username", data.Username)
	return token, nil
}

func isPasswordHard(pswrd string) error {
	if utf8.RuneCountInString(pswrd) < 8 {
		return ErrPasswordIsTooShort
	}
	return nil
}
