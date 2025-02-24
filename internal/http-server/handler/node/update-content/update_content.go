package updatecontent

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/http-server/api/validate"
	"main/internal/storage"
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

		err = nodeUpdater.UpdateNoteNodeContent(nodeId, req.Content)
		if errors.Is(err, storage.ErrNoteNodeNotFound) {
			log.Error("note node not found", "error", err)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(resperrors.ErrNodeDoesNotExist))

			return
		}
		if err != nil {
			log.Error("failed to update note node content", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToUpdateNodeContent))

			return
		}

		log.Info("note node content updated")

		render.JSON(w, r, resp.OK())
	}
}
