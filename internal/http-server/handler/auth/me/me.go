package me

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/models/user"
	"main/internal/storage"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	User user.User `json:"user"`
}

type UserGetter interface {
	GetUserById(id string) (user.User, error)
}

func New(log *slog.Logger, userGetter UserGetter, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		userId, ok := claims["user_id"].(string)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error(resperrors.ErrUserUnauthorized))
		}

		userFromDB, err := userGetter.GetUserById(userId)
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.Attr{Key: "user_id", Value: slog.StringValue(userId)})

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(resperrors.ErrUserDoesNotExist))

			return
		}
		if err != nil {
			log.Error("failed to get user", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError))

			return
		}

		render.JSON(w, r, Response{resp.OK(), userFromDB})
	}
}
