package services_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/undefeel/cloud-storage-backend/internal/lib/jwt"
	"github.com/undefeel/cloud-storage-backend/internal/repository"
	"github.com/undefeel/cloud-storage-backend/internal/services"
	mock_sevices "github.com/undefeel/cloud-storage-backend/internal/services/mocks"
	"testing"
)

func TestService_SignUp(t *testing.T) {
	cases := []struct {
		name        string
		input       *services.User
		expectedErr error
		authErr     error
		repoErr     error
	}{
		{
			name: "default",
			input: &services.User{
				Id:       uuid.New(),
				Login:    "testLogin",
				Password: "testPassword",
			},
		},
		{
			name: "empty",
			input: &services.User{
				Id:       uuid.New(),
				Login:    "",
				Password: "",
			},
			expectedErr: services.ErrIncorrectInput,
			repoErr:     repository.ErrValidation,
		},
		{
			name: "conflict",
			input: &services.User{
				Id:       uuid.New(),
				Login:    "conflictUser",
				Password: "password",
			},
			expectedErr: services.ErrUnique,
			repoErr:     repository.ErrUnique,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			auth := mock_sevices.NewAuthenticator(t)
			repo := mock_sevices.NewRepositoryProvider(t)
			repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*services.User")).Return(&services.User{}, c.repoErr).Once()
			if c.repoErr == nil {
				auth.On("NewPair", mock.AnythingOfType("uuid.UUID")).Return(jwt.TokenPair{}, c.authErr).Once()
			}

			s := services.New(repo, auth)

			_, err := s.SignUp(context.Background(), c.input)

			if err != nil {
				if !errors.Is(err, c.expectedErr) {
					t.Error("unexpected error")
				}
			}
		})
	}
}
