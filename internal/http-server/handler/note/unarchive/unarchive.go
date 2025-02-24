package unarchive

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/http-server/api/validate"
	"main/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type NoteUnarchiver interface {
	UnarchiveNote(id int) error
	validate.UserVerifier
}

func New(log *slog.Logger, noteUnarchiver NoteUnarchiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.unarchive.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		id, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		err = validate.VerifyUserNote(id, noteUnarchiver, w, r, log)
		if err != nil {
			return
		}

		err = noteUnarchiver.UnarchiveNote(id)
		if errors.Is(err, storage.ErrNoteNotFound) {
			log.Error("note not found", "error", err)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(resperrors.ErrNoteDoesNotExist))

			return
		}
		if err != nil {
			log.Error("failed to unarchive note", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToUnarchiveNote))

			return
		}

		log.Info("note unarchived", slog.Int("id", id))

		render.JSON(w, r, resp.OK())
	}
}
