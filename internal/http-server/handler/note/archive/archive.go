package archive

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
)

type NoteArchiver interface {
	ArchiveNote(id int) error
}

func New(log *slog.Logger, noteArchiver NoteArchiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.archive.New"

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

		err = noteArchiver.ArchiveNote(id)
		if errors.Is(err, storage.ErrNoteNotFound) {
			log.Error("note not found", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}
		if err != nil {
			log.Error("failed to archive note", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to archive note"))

			return
		}

		log.Info("note archived", slog.Int("id", id))

		render.JSON(w, r, resp.OK())
	}
}
