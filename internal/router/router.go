package router

import (
	"log/slog"
	"main/internal/config"
	"net/http"

	resp "main/internal/http-server/api/response"
	"main/internal/http-server/handler/node/add"
	"main/internal/http-server/handler/node/delete"
	updatecontent "main/internal/http-server/handler/node/update-content"
	"main/internal/http-server/handler/note/create"
	getnote "main/internal/http-server/handler/note/get-note"
	getusernotes "main/internal/http-server/handler/note/get-user-notes"
	updatefullnote "main/internal/http-server/handler/note/update-full-note"
	updatetitle "main/internal/http-server/handler/note/update-title"
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
	getnote.NoteGetter
	getusernotes.NotesGetter
	updatetitle.NoteTitleUpdater
	updatefullnote.NoteFUllUpdater

	add.NodeAdder
	delete.NodeDeleter
	updatecontent.NodeUpdater
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
	r.Get("/note/{id}", getnote.New(logger, storage))
	r.Get("/note/user/{id}", getusernotes.New(logger, storage))
	r.Patch("/note/{id}", updatetitle.New(logger, storage))
	r.Put("/note/{id}", updatefullnote.New(logger, storage))
	// TODO: add archive note patch handler
	// TODO: add unarchive note patch handler
	// TODO: add delete note delete handler

	r.Post("/node", add.New(logger, storage))
	r.Delete("/node", delete.New(logger, storage))
	r.Patch("/node/{id}", updatecontent.New(logger, storage))
}
