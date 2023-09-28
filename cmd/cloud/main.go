package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/undefeel/cloud-storage-backend/internal/config"
	"github.com/undefeel/cloud-storage-backend/internal/lib/jwt"
	"github.com/undefeel/cloud-storage-backend/internal/repository/postgres"
	"github.com/undefeel/cloud-storage-backend/internal/services"
	http_server "github.com/undefeel/cloud-storage-backend/internal/transport/http-server"
	dLog "log"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title           Cloud storage API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  Jwt

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cfg := config.MustLoad()

	ctx := context.Background()
	repo, err := postgres.New(ctx, cfg.Host, cfg.Port, cfg.User, cfg.DbName, cfg.Password)
	if err != nil {
		dLog.Fatalf("repository can not create: %s", err)
	}

	auth := jwt.New(cfg.AccessSecret, cfg.AccessExpire, cfg.RefreshSecret, cfg.RefreshExpire)

	service := services.New(repo, auth)

	log := setupLogger(cfg.Env)

	server := http_server.New(cfg.IdleTimeout, cfg.Timeout, cfg.BindAddr, log, service)
	server.MustServe()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
