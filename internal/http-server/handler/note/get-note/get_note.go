package getnote

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/http-server/api/validate"
	"main/internal/models/note"
	"main/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Note note.Note `json:"data"`
}

type NoteGetter interface {
	GetNoteById(id int) (note.Note, error)
	GetAllNotesNodes(noteId int) ([]note.NoteNode, error)
	validate.UserVerifier
}

func New(log *slog.Logger, noteGetter NoteGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		id, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		noteFromDB, err := noteGetter.GetNoteById(id)
		if errors.Is(err, storage.ErrNoteNotFound) {
			log.Error("note not found", "error", err)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(resperrors.ErrNoteDoesNotExist))

			return
		}
		if err != nil {
			log.Error("failed to get note", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToGetNote))

			return
		}

		err = validate.VerifyUserNote(noteFromDB.Id, noteGetter, w, r, log)
		if err != nil {
			return
		}

		nodes, err := noteGetter.GetAllNotesNodes(noteFromDB.Id)
		if errors.Is(err, storage.ErrNoteNodeNotFound) {
			log.Error("note nodes not found", "error", err)

			nodes = []note.NoteNode{}
		}
		if err != nil {
			log.Error("failed to get note nodes", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToGetNoteNodes))

			return
		}

		for i, n := range nodes {
			if n.ContentType == note.ContentTypeImage {
				nodes[i].Content = ""
			}
		}

		noteFromDB.Nodes = nodes

		log.Info("note got", slog.Int("id", noteFromDB.Id))

		render.JSON(w, r, Response{resp.OK(), noteFromDB})
	}
}
