package delete

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

type NodeDeleter interface {
	DeleteNoteNode(int) error
	validate.UserVerifier
}

func New(log *slog.Logger, nodeDeleter NodeDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.node.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		id, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		err = validate.VerifyUserNoteNode(id, nodeDeleter, w, r, log)
		if err != nil {
			return
		}

		err = nodeDeleter.DeleteNoteNode(id)
		if errors.Is(err, storage.ErrNoteNodeNotFound) {
			log.Error("not found note node", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(resperrors.ErrNodeDoesNotExist))

			return
		}
		if err != nil {
			log.Error("failed to delete note node", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrFailedToDeleteNode))

			return
		}

		log.Info("note node deleted")

		render.JSON(w, r, resp.OK())
	}
}
