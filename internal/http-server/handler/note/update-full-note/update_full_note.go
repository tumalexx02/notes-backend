package updatefullnote

import (
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/models/note"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Note note.Note `note:"required"`
}

type NoteFUllUpdater interface {
	UpdateFullNote(id int, note note.Note) (int, error)
}

func New(log *slog.Logger, noteUpdater NoteFUllUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.updatefullnote.New"

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

		log.Debug("got note", slog.Any("note", req.Note))

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 {
			log.Error("invalid 'id' param", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("invalid 'id' param"))

			return
		}

		rows, err := noteUpdater.UpdateFullNote(id, req.Note)

		if err != nil {
			log.Error("failed to update note", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to update note"))

			return
		}

		log.Info("note updated", slog.Int("id", id), slog.Int("rows_affected", rows))

		render.JSON(w, r, resp.OK())
	}
}
