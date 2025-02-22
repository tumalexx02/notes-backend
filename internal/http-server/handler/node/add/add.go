package add

import (
	"log/slog"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/http-server/api/validate"
	"main/internal/models/note"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
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
	validate.UserVerifier
}

func New(log *slog.Logger, noteAdder NodeAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.node.add.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := validate.DecodeAndValidateRequestJson(&req, w, r, log); err != nil {
			return
		}

		noteId := req.NoteId

		err := validate.VerifyUserNote(noteId, noteAdder, w, r, log)
		if err != nil {
			return
		}

		contentType := req.ContentType
		content := req.Content

		if contentType == note.ContentTypeImage {
			content = ""
		}

		id, err := noteAdder.AddNoteNode(noteId, contentType, content)
		if err != nil {
			log.Error("failed to add note node", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToAddNoteNode))

			return
		}

		log.Info("node added", slog.Int("node_id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Id:       id,
		})
	}
}
