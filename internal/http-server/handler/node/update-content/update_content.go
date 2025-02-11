package updatecontent

import (
	"log/slog"
	resp "main/internal/http-server/api/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Content string `json:"content"`
}

type NodeUpdater interface {
	UpdateNoteNodeContent(id int, content string) error
}

func New(log *slog.Logger, nodeUpdater NodeUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.node.updatecontent.New"

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

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 {
			log.Error("invalid 'id' param", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("invalid 'id' param"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if err := nodeUpdater.UpdateNoteNodeContent(id, req.Content); err != nil {
			log.Error("failed to update note node content", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to update note node content"))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Info("note node content updated")

		render.JSON(w, r, resp.OK())
	}
}
