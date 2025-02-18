package router

import (
	"log/slog"
	"main/internal/config"
	"net/http"

	resp "main/internal/http-server/api/response"
	loggerMiddleware "main/internal/http-server/middleware/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type Router struct {
	*chi.Mux
	jwtauth *jwtauth.JWTAuth
}

type Storage interface {
	Noter
	NoteNoder
	Authorizer
}

func New(cfg *config.Config, log *slog.Logger) *Router {
	// init chi router
	router := chi.NewRouter()

	// add middlewares
	router.Use(middleware.RequestID)
	router.Use(loggerMiddleware.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	return &Router{
		router,
		generateAuthToken(cfg),
	}
}

func (r *Router) InitRoutes(storage Storage, logger *slog.Logger, cfg *config.Config) {
	// health check route
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, resp.OK())
	})

	r.InitAuthRoutes(storage, logger, cfg)
	r.InitNotesRoutes(storage, logger, cfg)
	r.InitNoteNodesRoutes(storage, logger, cfg)
}
