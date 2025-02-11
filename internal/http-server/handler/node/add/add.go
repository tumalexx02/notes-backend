package add

import (
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/models/note"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	NoteId      int    `json:"note_id" validate:"required"`
	ContentType string `json:"content_type" validate:"required,custom_url"`
	Content     string `json:"content"`
}

type Response struct {
	resp.Response
	Id int `json:"node_id"`
}

type NodeAdder interface {
	AddNoteNode(noteId int, contentType string, content string) (int, error)
}

func New(log *slog.Logger, noteAdder NodeAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.node.add.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to decode request body"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		val := validator.New()
		if err := val.RegisterValidation("custom_url", categoryValidator); err != nil {
			log.Error("validator init error", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			return
		}

		if err := val.Struct(req); err != nil {
			log.Error("invalid request body", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("invalid request body"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		noteId := req.NoteId
		contentType := req.ContentType
		content := req.Content

		id, err := noteAdder.AddNoteNode(noteId, contentType, content)
		if err != nil {
			log.Error("failed to add note node", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to add note node"))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Info("node added", slog.Int("node_id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Id:       id,
		})
	}
}

func categoryValidator(fl validator.FieldLevel) bool {
	category := fl.Field().String()
	switch category {
	case note.ContentTypeImage, note.ContentTypeText, note.ContentTypeList:
		return true
	default:
		return false
	}
}
