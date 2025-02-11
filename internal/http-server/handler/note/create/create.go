package create

import (
	"log/slog"
	resp "main/internal/http-server/api/response"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Title  string `json:"title" validate:"required,max=31"`
	UserId string `json:"user_id" validate:"required"`
}

type Response struct {
	resp.Response
	Id uint `json:"note_id"`
}

type NoteCreator interface {
	CreateNote(noteTitle string, userId string) (uint, error)
}

func New(log *slog.Logger, storage NoteCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.create.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to decode request body"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request body", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("invalid request body"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		title := req.Title
		userId := req.UserId

		id, err := storage.CreateNote(title, userId)
		if err != nil {
			log.Error("failed to create note", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to create note"))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Info("note created", slog.Int("note_id", int(id)))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Id:       id,
		})
	}
}
