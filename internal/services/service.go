package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/undefeel/cloud-storage-backend/internal/lib/jwt"
	"github.com/undefeel/cloud-storage-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo RepositoryProvider
	auth Authenticator
}

//go:generate go run github.com/vektra/mockery/v2@v2.33.3 --all
type RepositoryProvider interface {
	CreateUser(ctx context.Context, us *User) (*User, error)
	FindUserByLogin(ctx context.Context, us *User) (*User, error)
}

type Authenticator interface {
	NewPair(userID uuid.UUID) (jwt.TokenPair, error)
}

func New(repo RepositoryProvider, auth Authenticator) Service {
	return Service{
		repo: repo,
		auth: auth,
	}
}

func (s Service) SignUp(ctx context.Context, user *User) (jwt.TokenPair, error) {
	const op = "services.service.SignUp"
	hPassword, err := genHash(user.Password)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return jwt.TokenPair{}, fmt.Errorf("%w: %w", ErrIncorrectInput, err)
	}
	user.Password = hPassword
	us, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return jwt.TokenPair{}, switchRepoErrors(err)
	}

	pair, err := s.auth.NewPair(us.Id)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return jwt.TokenPair{}, switchJwtErrors(err)
	}
	return pair, nil
}

// SignIn create user in repository, compare passwords and return jwt tokens pair
//
// Errors types: ErrIncorrectInput, ErrNotFound
func (s Service) SignIn(ctx context.Context, user *User) (jwt.TokenPair, error) {
	const op = "services.service.SignIn"
	us, err := s.repo.FindUserByLogin(ctx, user)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return jwt.TokenPair{}, switchRepoErrors(err)
	}
	err = checkHash(user.Password, us.Password)
	if err != nil {
		return jwt.TokenPair{}, fmt.Errorf("%s: %w %w", op, ErrIncorrectInput, err)
	}

	pair, err := s.auth.NewPair(us.Id)
	if err != nil {
		return jwt.TokenPair{}, fmt.Errorf("%s: %w", op, err)
	}

	return pair, nil
}

// Max len password must be not longer than 72 bytes (72 symbols ascii)
func genHash(password string) (string, error) {
	const op = "services.service.genHash"
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return string(bytes), nil
}

func checkHash(password string, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

func switchRepoErrors(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	case errors.Is(err, repository.ErrUnique):
		return fmt.Errorf("%w: %w", ErrUnique, err)
	case errors.Is(err, repository.ErrValidation):
		return fmt.Errorf("%w: %w", ErrIncorrectInput, err)
	default:
		return err
	}
}

func switchJwtErrors(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return fmt.Errorf("%w: %w", ErrTokenExpired, err)
	case errors.Is(err, jwt.ErrTokenInvalid):
		return fmt.Errorf("%w: %w", ErrTokenInvalid, err)
	default:
		return err
	}
}
