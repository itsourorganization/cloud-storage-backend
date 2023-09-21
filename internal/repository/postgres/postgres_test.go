package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/undefeel/cloud-storage-backend/internal/repository"
	"github.com/undefeel/cloud-storage-backend/internal/repository/postgres/ent"
	"github.com/undefeel/cloud-storage-backend/internal/repository/postgres/ent/enttest"
	"github.com/undefeel/cloud-storage-backend/internal/services"
	"os"
	"testing"
)

type cfgForTest struct {
	Host     string `env:"DATABASE_HOST" env-required:"true"`
	Port     string `env:"DATABASE_PORT" env-required:"true"`
	User     string `env:"DATABASE_USER" env-required:"true"`
	DbName   string `env:"DATABASE_NAME" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
}

func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}

func makeCl(t *testing.T) *ent.Client {
	var cfg cfgForTest
	err := cleanenv.ReadEnv(&cfg)
	require.NoError(t, err)
	dbOpts := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DbName,
		cfg.Password)
	return enttest.Open(t, "postgres", dbOpts)

}

func TestRepository_CreateUser(t *testing.T) {
	cases := []struct {
		name        string
		input       *services.User
		expectedErr error
	}{
		{
			name: "default",
			input: &services.User{
				Id:       uuid.New(),
				Login:    "testUser",
				Password: "testPassword",
			},
			expectedErr: nil,
		},
		{
			name: "unique",
			input: &services.User{
				Id:       uuid.New(),
				Login:    "testUser",
				Password: "testPassword",
			},
			expectedErr: repository.ErrUnique,
		},
		{
			name: "empty",
			input: &services.User{
				Id:       uuid.New(),
				Login:    "",
				Password: "",
			},
			expectedErr: repository.ErrValidation,
		},
	}
	cl := makeCl(t)
	defer cl.User.Delete().Exec(context.Background())
	repo := Repository{cl: cl}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			_, err := repo.CreateUser(context.Background(), cs.input)
			if err != nil {
				if errors.Is(err, cs.expectedErr) {
					return
				} else {
					t.Fatal("bad type error", err)
				}
			}
		})
	}
}
