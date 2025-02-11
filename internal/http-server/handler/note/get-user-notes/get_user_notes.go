package getusernotes

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/models/note"
	"main/internal/storage"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	resp.Response
	Notes []note.NotePreview `json:"data"`
}

type NotesGetter interface {
	GetUserNotes(userId string) ([]note.NotePreview, error)
}

func New(log *slog.Logger, notesGetter NotesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if err := validator.New().Var(id, "uuid"); err != nil {
			log.Error("invalid 'id' param", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("invalid 'id' param"))

			return
		}

		notes, err := notesGetter.GetUserNotes(id)
		if errors.Is(err, storage.ErrNoteNotFound) {
			log.Error("note not found", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}
		if err != nil {
			log.Error("failed to get notes", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to get notes"))

			return
		}

		render.JSON(w, r, Response{
			resp.OK(),
			notes,
		})
	}
}
