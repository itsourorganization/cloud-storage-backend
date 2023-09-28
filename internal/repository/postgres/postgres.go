package postgres

import (
	"context"
	"fmt"
	"github.com/undefeel/cloud-storage-backend/internal/repository"
	"github.com/undefeel/cloud-storage-backend/internal/repository/postgres/ent"
	"github.com/undefeel/cloud-storage-backend/internal/repository/postgres/ent/migrate"
	"github.com/undefeel/cloud-storage-backend/internal/repository/postgres/ent/user"
	"github.com/undefeel/cloud-storage-backend/internal/services"
)

type Repository struct {
	cl *ent.Client
}

func New(ctx context.Context, host, port, user, dbName, password string) (Repository, error) {
	dbOps := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbName, password)
	cl, err := ent.Open("postgres", dbOps)
	if err != nil {
		return Repository{}, err
	}
	if err := cl.Schema.Create(ctx, migrate.WithDropColumn(true), migrate.WithDropIndex(true)); err != nil {
		return Repository{}, err
	}

	return Repository{cl: cl}, nil
}

func convertEntUserToService(user *ent.User) *services.User {
	return &services.User{
		Id:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	}
}

func (r Repository) CreateUser(ctx context.Context, us *services.User) (*services.User, error) {
	const op = "repository.postgres.CreateUser"
	dbUser, err := r.cl.User.Create().SetLogin(us.Login).SetPassword(us.Password).SetID(us.Id).Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fmt.Errorf("%s: %w %w", op, repository.ErrUnique, err)
		}
		if ent.IsValidationError(err) {
			return nil, fmt.Errorf("%s: %w %w", op, repository.ErrValidation, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return convertEntUserToService(dbUser), nil
}

func (r Repository) FindUserByLogin(ctx context.Context, us *services.User) (*services.User, error) {
	const op = "repository.postgres.FingUserByLogin"
	dbUs, err := r.cl.User.Query().Where(user.Login(us.Login)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("%s: %w %w", op, repository.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return convertEntUserToService(dbUs), nil
}
