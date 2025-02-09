package router

import (
	"log/slog"
	"main/internal/config"
	"net/http"

	loggerMiddleware "main/internal/http-server/middleware/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	*chi.Mux
}

type Storage interface {
}

func New(cfg *config.Config, log *slog.Logger) *Router {
	if log == nil {
		panic("logger is nil")
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(loggerMiddleware.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	return &Router{
		router,
	}
}

func (r *Router) InitRoutes(storage Storage, logger *slog.Logger, cfg *config.Config) {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
}
