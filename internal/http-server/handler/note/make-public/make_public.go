package makepublic

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/http-server/api/validate"
	"main/internal/storage"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	PublicId string `json:"public_id"`
}

type PublicNoteMaker interface {
	MakeNotePublic(noteId int) (string, error)
	validate.UserVerifier
}

func New(log *slog.Logger, publicNoteMaker PublicNoteMaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.makepublic.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		id, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		err = validate.VerifyUserNote(id, publicNoteMaker, w, r, log)
		if err != nil {
			return
		}

		publicId, err := publicNoteMaker.MakeNotePublic(id)
		if errors.Is(err, storage.ErrNoteNotFound) {
			log.Error("note not found", "error", err)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(resperrors.ErrNoteDoesNotExist))

			return
		}
		if err != nil {
			log.Error("failed to make note public", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToMakeNotePublic))

			return
		}

		render.JSON(w, r, Response{resp.OK(), publicId})
	}
}
