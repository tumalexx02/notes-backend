package authenticator

import (
	"errors"
	"fmt"
	"log/slog"
	resp "main/internal/http-server/api/response"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

func Authenticator(ja *jwtauth.JWTAuth, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/authenticator"),
		)

		log.Info("authenticator middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				log.Error("auth context error", "error", fmt.Sprintf("%+v", err))

				w.WriteHeader(http.StatusUnauthorized)

				if errors.Is(err, jwtauth.ErrExpired) {
					render.JSON(w, r, resp.Error("token expired"))
				} else {
					render.JSON(w, r, resp.Error("invalid token"))
				}

				return
			}

			if token == nil {
				log.Error("token not found")

				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("invalid token"))

				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
