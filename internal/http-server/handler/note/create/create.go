package create

import (
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/http-server/api/validate"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type Request struct {
	Title string `json:"title" validate:"required,max=31"`
}

type Response struct {
	resp.Response
	Id int `json:"note_id"`
}

type NoteCreator interface {
	CreateNote(noteTitle string, userId string) (int, error)
}

func New(log *slog.Logger, noteCreator NoteCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.create.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := validate.DecodeAndValidateRequestJson(&req, w, r, log); err != nil {
			return
		}

		title := req.Title

		_, claims, _ := jwtauth.FromContext(r.Context())

		userId, _ := claims["user_id"].(string)

		id, err := noteCreator.CreateNote(title, userId)
		if err != nil {
			log.Error("failed to create note", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to create note"))

			return
		}

		log.Info("note created", slog.Int("note_id", int(id)))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Id:       id,
		})
	}
}
