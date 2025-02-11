package router

import (
	"log/slog"
	"main/internal/config"
	"net/http"

	resp "main/internal/http-server/api/response"
	"main/internal/http-server/handler/node/add"
	"main/internal/http-server/handler/node/delete"
	"main/internal/http-server/handler/note/create"
	loggerMiddleware "main/internal/http-server/middleware/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Router struct {
	*chi.Mux
}

type Storage interface {
	create.NoteCreator
	add.NoteAdder
	delete.NodeDeleter
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
		render.JSON(w, r, resp.OK())
	})

	r.Post("/note", create.New(logger, storage))
	r.Post("/node", add.New(logger, storage))
	r.Delete("/node", delete.New(logger, storage))
}
