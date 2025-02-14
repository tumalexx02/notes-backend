package updateorder

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
	OldOrder int `json:"old_order" validate:"gte=0"`
	NewOrder int `json:"new_order" validate:"gte=0"`
}

type NoteOrderUpdater interface {
	UpdateNoteNodeOrder(noteId int, oldOrder int, newOrder int) error
	validate.UserVerifier
}

func New(log *slog.Logger, noteUpdater NoteOrderUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.updateorder.New"

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

		err = validate.VerifyUserNote(id, noteUpdater, w, r, log)
		if err != nil {
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
