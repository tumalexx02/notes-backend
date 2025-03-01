package uploadimage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log/slog"
	"main/internal/config"
	"main/internal/http-server/api/response"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/http-server/api/validate"
	"main/internal/models/note"
	"main/internal/storage"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type ImageUploader interface {
	UpdateNoteNodeContent(id int, content string) error
	GetNodeById(id int) (note.NoteNode, error)
	validate.UserVerifier
}

func New(cfg *config.Config, log *slog.Logger, imageUploader ImageUploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.node.uploadimage.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		id, err := validate.GetIntURLParam("id", w, r, log)
		if err != nil {
			return
		}

		err = validate.VerifyUserNoteNode(id, imageUploader, w, r, log)
		if err != nil {
			return
		}

		noteFromDB, err := imageUploader.GetNodeById(id)
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

		if noteFromDB.ContentType != note.ContentTypeImage {
			log.Info("invalid content type", slog.Any("content-type", noteFromDB.ContentType))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(resperrors.ErrInvalidContentType))

			return
		}

		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error("failed to parse multipart form", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError))

			return
		}

		imageFile, format, err := r.FormFile("image")
		if err != nil {
			log.Error("failed to get image", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError))

			return
		}
		defer imageFile.Close()

		imageType := format.Header.Get("Content-Type")

		if imageType != "image/jpeg" && imageType != "image/png" {
			log.Error("invalid image format", slog.String("format", format.Header.Get("Content-Type")))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(resperrors.ErrInvalidImageFormat))

			return
		}

		exp := ".jpg"
		if imageType == "image/png" {
			exp = ".png"
		}

		fileId := uuid.New().String()
		fileName := fileId + exp
		dirName := filepath.Join(cfg.Image.ImagesDir, HashNoteId(noteFromDB.NoteId, cfg.Image.ImageSalt), strconv.Itoa(noteFromDB.Id))
		imagePath := filepath.Join(dirName, fileName)

		if err := os.RemoveAll(dirName); err != nil {
			log.Error("failed to remove existing directory", slog.String("error", err.Error()))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error(resperrors.ErrInternalServerError))

			return
		}

		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			log.Error("failed to create image dir", slog.String("error", err.Error()))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error(resperrors.ErrInternalServerError))

			return
		}

		outFile, err := os.Create(imagePath)
		if err != nil {
			log.Error("failed to create image file", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error(resperrors.ErrInternalServerError))
			return
		}
		defer outFile.Close()

		err = compressImage(imageFile, imageType, outFile.Name(), cfg.Image.MaxWidth)
		if err != nil {
			log.Error("failed to compress image", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error(resperrors.ErrInternalServerError))
			return
		}

		err = imageUploader.UpdateNoteNodeContent(id, imagePath)
		if err != nil {
			log.Error("failed to update note node content", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error(resperrors.ErrInternalServerError))
			return
		}

		file, err := os.Open(imagePath)
		if err != nil {
			log.Error("failed to open image file", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error(resperrors.ErrInternalServerError))
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Error("failed to get file info", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error(resperrors.ErrInternalServerError))
			return
		}

		w.Header().Set("Content-Disposition", "inline; filename="+fileName)
		w.Header().Set("Content-Type", "image/"+imageType)

		http.ServeContent(w, r, fileName, fileInfo.ModTime(), file)

		render.JSON(w, r, resp.OK())
	}
}

func HashNoteId(noteId int, salt string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(strconv.Itoa(noteId)))
	return hex.EncodeToString(h.Sum(nil))
}

func compressImage(file multipart.File, format string, savePath string, maxWidth uint) error {
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Масштабирование изображения до maxWidth
	resizedImg := resize.Resize(maxWidth, 0, img, resize.Lanczos3)

	outFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if format == "image/jpeg" {
		return jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 60})
	} else if format == "image/png" {
		return png.Encode(outFile, resizedImg)
	}
	return nil
}
