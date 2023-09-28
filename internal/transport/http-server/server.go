package http_server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/undefeel/cloud-storage-backend/internal/services"
	"github.com/undefeel/cloud-storage-backend/internal/transport/http-server/handlers/signIn"
	signup "github.com/undefeel/cloud-storage-backend/internal/transport/http-server/handlers/signUp"
	"github.com/undefeel/cloud-storage-backend/internal/transport/http-server/middleware/logger"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	s       *http.Server
	log     *slog.Logger
	service services.Service
}

func New(idleTimeout time.Duration, timeout time.Duration, addr string, log *slog.Logger, u services.Service) Server {
	s := &http.Server{
		ReadTimeout:  timeout,
		IdleTimeout:  idleTimeout,
		WriteTimeout: timeout,
		Addr:         addr,
	}
	return Server{s: s, log: log, service: u}
}

func (s Server) MustServe() {
	r := chi.NewRouter()
	r.Use(logger.New(s.log))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/signUp", signup.New(s.log, s.service))
		r.Post("/signIn", signIn.New(s.log, s.service))
	})

	s.s.Handler = r
	err := s.s.ListenAndServe()
	if err != nil {
		log.Fatalf("server can not run: %s", err)
	}
}
