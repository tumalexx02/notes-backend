package router

import (
	"log/slog"
	"main/internal/config"
	"main/internal/http-server/handler/note/archive"
	"main/internal/http-server/handler/note/create"
	deleteNote "main/internal/http-server/handler/note/delete"
	getnote "main/internal/http-server/handler/note/get-note"
	getusernotes "main/internal/http-server/handler/note/get-user-notes"
	"main/internal/http-server/handler/note/unarchive"
	updatefullnote "main/internal/http-server/handler/note/update-full-note"
	updateorder "main/internal/http-server/handler/note/update-order"
	updatetitle "main/internal/http-server/handler/note/update-title"
	"main/internal/http-server/middleware/authenticator"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

type Noter interface {
	create.NoteCreator
	getnote.NoteGetter
	getusernotes.NotesGetter
	updatetitle.NoteTitleUpdater
	updatefullnote.NoteFUllUpdater
	archive.NoteArchiver
	unarchive.NoteUnarchiver
	deleteNote.NoteDeleter
	updateorder.NoteOrderUpdater
}

func (r *Router) InitNotesRoutes(storage Storage, logger *slog.Logger, cfg *config.Config) {
	// note routes
	r.Route("/note", func(noteRouter chi.Router) {
		noteRouter.Use(jwtauth.Verifier(r.jwtauth))
		noteRouter.Use(authenticator.Authenticator(r.jwtauth, logger))

		// create
		noteRouter.Post("/create", create.New(logger, storage))

		// read
		noteRouter.Get("/{id}", getnote.New(logger, storage))
		noteRouter.Get("/list", getusernotes.New(logger, storage))

		// update
		noteRouter.Put("/{id}", updatefullnote.New(logger, storage))
		noteRouter.Patch("/{id}", updatetitle.New(logger, storage))
		noteRouter.Patch("/{id}/order", updateorder.New(logger, storage))

		// delete (and archive)
		noteRouter.Patch("/{id}/archive", archive.New(logger, storage))
		noteRouter.Patch("/{id}/unarchive", unarchive.New(logger, storage))
		noteRouter.Delete("/{id}", deleteNote.New(logger, storage))
	})
}
