package updatetitle

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/http-server/api/validate"
	"main/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Title string `json:"title" validate:"required,max=31"`
}

type NoteTitleUpdater interface {
	UpdateNoteTitle(id int, title string) error
	validate.UserVerifier
}

func New(log *slog.Logger, noteTitleUpdater NoteTitleUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.updatetitle.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := validate.DecodeAndValidateRequestJson(&req, w, r, log); err != nil {
			return
		}

		id, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		err = validate.VerifyUser(id, noteTitleUpdater, w, r, log)
		if err != nil {
			return
		}

		title := req.Title

		err = noteTitleUpdater.UpdateNoteTitle(id, title)
		if errors.Is(err, storage.ErrNoteNotFound) {
			log.Error("note not found", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}
		if err != nil {
			log.Error("failed to update note title", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to update note title"))

			return
		}

		log.Info("note title updated", slog.Int("id", id))

		render.JSON(w, r, resp.OK())
	}
}
