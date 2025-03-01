package makeprivate

import (
	"errors"
	"log/slog"
	"main/internal/http-server/api/validate"
	"main/internal/storage"
	"net/http"

	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type PrivateNoteMaker interface {
	MakeNotePrivate(id int) error
	validate.UserVerifier
}

func New(log *slog.Logger, privateNoteMaker PrivateNoteMaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.makeprivate.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		id, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		err = validate.VerifyUserNote(id, privateNoteMaker, w, r, log)
		if err != nil {
			return
		}

		err = privateNoteMaker.MakeNotePrivate(id)
		if errors.Is(err, storage.ErrNoteNotFound) {
			log.Error("note not found", "error", err)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(resperrors.ErrNoteDoesNotExist))

			return
		}
		if err != nil {
			log.Error("failed to make note private", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToMakeNotePrivate))

			return
		}

		render.JSON(w, r, resp.OK())
	}
}
