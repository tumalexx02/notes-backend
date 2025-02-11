package updatetitle

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/storage"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Title string `json:"title" validate:"required,max=31"`
}

type NoteTitleUpdater interface {
	UpdateNoteTitle(id int, title string) error
}

func New(log *slog.Logger, noteTitleUpdater NoteTitleUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.updatetitle.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to decode request body"))

			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request body", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("invalid request body"))

			return
		}

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 {
			log.Error("invalid 'id' param", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("invalid 'id' param"))

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

		log.Info("note title updated")

		render.JSON(w, r, resp.OK())
	}
}
