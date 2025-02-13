package me

import (
	"log/slog"
	resp "main/internal/http-server/api/response"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	UserId string `json:"user_id"`
}

func New(log *slog.Logger, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		userId, ok := claims["user_id"].(string)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("user unauthorized"))
		}

		render.JSON(w, r, Response{resp.OK(), userId})
	}
}
