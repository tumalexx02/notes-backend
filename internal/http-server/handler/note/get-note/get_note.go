package getnote

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/models/note"
	"main/internal/storage"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Note note.Note `json:"data"`
}

type NoteGetter interface {
	GetNoteById(id int) (note.Note, error)
}

func New(log *slog.Logger, noteGetter NoteGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.get.New"

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

		note, err := noteGetter.GetNoteById(id)
		if errors.Is(err, storage.ErrNoteNotFound) {
			log.Error("note not found", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}
		if err != nil {
			log.Error("failed to get note", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to get note"))

			return
		}

		log.Info("note got", slog.Int("id", note.Id))

		render.JSON(w, r, Response{resp.OK(), note})
	}
}
