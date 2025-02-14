package updatecontent

import (
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/http-server/api/validate"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Content string `json:"content"`
}

type NodeUpdater interface {
	UpdateNoteNodeContent(id int, content string) error
	validate.UserVerifier
}

func New(log *slog.Logger, nodeUpdater NodeUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.node.updatecontent.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := validate.DecodeRequestJson(&req, w, r, log); err != nil {
			return
		}

		nodeId, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		err = validate.VerifyUserNoteNode(nodeId, nodeUpdater, w, r, log)
		if err != nil {
			return
		}

		if err := nodeUpdater.UpdateNoteNodeContent(nodeId, req.Content); err != nil {
			log.Error("failed to update note node content", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to update note node content"))

			return
		}

		log.Info("note node content updated")

		render.JSON(w, r, resp.OK())
	}
}
