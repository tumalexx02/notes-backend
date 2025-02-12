package updateorder

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/storage"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	OldOrder int `json:"old_order" validate:"gte=0"`
	NewOrder int `json:"new_order" validate:"gte=0"`
}

type NoteOrderUpdater interface {
	UpdateNoteNodeOrder(noteId int, oldOrder int, newOrder int) error
}

func New(log *slog.Logger, noteUpdater NoteOrderUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.updateorder.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 {
			log.Error("invalid 'id' param", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("invalid 'id' param"))

			return
		}

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

		err = noteUpdater.UpdateNoteNodeOrder(id, req.OldOrder, req.NewOrder)
		if errors.Is(err, storage.ErrNoteNodeNotFound) {
			log.Error("note node not found", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}
		if err != nil {
			log.Error("failed to update note node order", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to update note node order"))

			return
		}

		log.Info("note node order updated", slog.Int("noteId", id), slog.Int("oldOrder", req.OldOrder), slog.Int("newOrder", req.NewOrder))

		render.JSON(w, r, resp.OK())
	}
}
