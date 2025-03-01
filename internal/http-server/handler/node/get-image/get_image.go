package getimage

import (
	"errors"
	"log/slog"
	"main/internal/http-server/api/validate"
	"main/internal/models/note"
	"main/internal/storage"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type ImageGetter interface {
	GetNodeById(id int) (note.NoteNode, error)
	validate.UserVerifier
}

func New(log *slog.Logger, imageGetter ImageGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.node.getimage.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		id, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		err = validate.VerifyUserNoteNode(id, imageGetter, w, r, log)
		if err != nil {
			return
		}

		node, err := imageGetter.GetNodeById(id)
		if errors.Is(err, storage.ErrNoteNodeNotFound) {
			log.Error("note node not found", "error", err)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(resperrors.ErrNodeDoesNotExist))

			return
		}
		if err != nil {
			log.Error("failed to get note id", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError))

			return
		}

		if node.ContentType != note.ContentTypeImage {
			log.Error("node is not image", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(resperrors.ErrNodeIsNotImage))

			return
		}

		file, err := os.Open(node.Content)
		if err != nil {
			log.Error("failed to open image file", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError))
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Error("failed to get file info", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError))
			return
		}

		fileName := filepath.Base(file.Name())
		ext := filepath.Ext(fileName)
		mimeType := mime.TypeByExtension(ext)

		w.Header().Set("Content-Disposition", "inline; filename="+fileName)
		w.Header().Set("Content-Type", mimeType)

		http.ServeContent(w, r, fileName, fileInfo.ModTime(), file)
	}
}
