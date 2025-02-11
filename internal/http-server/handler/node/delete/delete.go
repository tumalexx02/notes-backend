package delete

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Id int `json:"id" validate:"required"`
}

type NodeDeleter interface {
	DeleteNoteNode(id int) error
}

func New(log *slog.Logger, nodeDeleter NodeDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.node.delete.New"

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

		err := nodeDeleter.DeleteNoteNode(req.Id)
		if errors.Is(err, storage.ErrNoteNodeNotFound) {
			log.Error("not found note node", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}
		if err != nil {
			log.Error("failed to delete note node", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to delete note node"))

			return
		}

		log.Info("note node deleted")

		render.JSON(w, r, resp.OK())
	}
}
