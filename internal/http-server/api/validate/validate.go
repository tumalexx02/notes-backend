package validate

import (
	"fmt"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/models/note"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UserVerifier interface {
	IsUserNoteOwner(userId string, noteId int) (bool, error)
}

func DecodeRequestJson[T any](dest *T, w http.ResponseWriter, r *http.Request, log *slog.Logger) error {
	if err := render.DecodeJSON(r.Body, dest); err != nil {
		log.Error("failed to decode request body", "error", err)

		render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError.Error()))

		return err
	}

	return nil
}

func DecodeAndValidateRequestJson[T any](dest *T, w http.ResponseWriter, r *http.Request, log *slog.Logger) error {
	err := DecodeRequestJson[T](dest, w, r, log)
	if err != nil {
		return err
	}

	val := validator.New()
	if err := val.RegisterValidation("custom_url", categoryValidator); err != nil {
		log.Error("validator init error", "error", err)

		render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError.Error()))

		return err
	}

	if err := val.Struct(*dest); err != nil {
		log.Error("invalid request body", "error", err)

		render.JSON(w, r, resp.Error(resperrors.ErrInvalidRequestBody.Error()))

		return err
	}

	return nil
}

func GetIntURLParam(paramName string, w http.ResponseWriter, r *http.Request, log *slog.Logger) (int, error) {
	strParam := chi.URLParam(r, paramName)

	intParam, err := strconv.Atoi(strParam)
	if err != nil || intParam < 0 {
		paramError := fmt.Errorf("invalid '%s' param", paramName)

		log.Error(paramError.Error(), "error", err)

		render.JSON(w, r, resp.Error(paramError.Error()))

		return 0, err
	}

	return intParam, nil
}

func VerifyUser(id int, userVerifier UserVerifier, w http.ResponseWriter, r *http.Request, log *slog.Logger) error {
	_, claims, _ := jwtauth.FromContext(r.Context())

	userId, _ := claims["user_id"].(string)

	isOwner, err := userVerifier.IsUserNoteOwner(userId, id)
	if err != nil {
		log.Error("failed to check note owner", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

		render.JSON(w, r, resp.Error("failed to check note owner"))

		return err
	}

	if !isOwner {
		log.Error("user is not note owner", "error", resperrors.ErrUserNotOwner, "user_id", userId, "note_id", id)

		render.JSON(w, r, resp.Error(resperrors.ErrUserNotOwner.Error()))

		return resperrors.ErrUserNotOwner
	}

	return nil
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
