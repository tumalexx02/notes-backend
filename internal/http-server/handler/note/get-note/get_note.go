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

		// TODO: separate get note and get nodes inside
		// TODO: on getting nodes add image node logic
		note, err := noteGetter.GetNoteById(id)
		if errors.Is(err, storage.ErrNoteNotFound) {
			log.Error("note not found", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(resperrors.ErrNoteDoesNotExist))

			return
		}
		if err != nil {
			log.Error("failed to get note", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToGetNote))

			return
		}

		err = validate.VerifyUserNote(note.Id, noteGetter, w, r, log)
		if err != nil {
			return
		}

		log.Info("note got", slog.Int("id", note.Id))

		render.JSON(w, r, Response{resp.OK(), note})
	}
}
