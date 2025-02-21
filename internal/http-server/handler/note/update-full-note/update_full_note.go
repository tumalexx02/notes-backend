package updatefullnote

import (
	"log/slog"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/http-server/api/validate"
	"main/internal/models/note"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Note note.Note `note:"required"`
}

type NoteFUllUpdater interface {
	UpdateFullNote(id int, note note.Note) (int, error)
	validate.UserVerifier
}

func New(log *slog.Logger, noteUpdater NoteFUllUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.updatefullnote.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := validate.DecodeRequestJson(&req, w, r, log); err != nil {
			return
		}

		id, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		err = validate.VerifyUserNote(id, noteUpdater, w, r, log)
		if err != nil {
			return
		}

		rows, err := noteUpdater.UpdateFullNote(id, req.Note)
		if err != nil {
			log.Error("failed to update note", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToUpdateFullNote))

			return
		}

		log.Info("note updated", slog.Int("id", id), slog.Int("rows_affected", rows))

		render.JSON(w, r, resp.OK())
	}
}
