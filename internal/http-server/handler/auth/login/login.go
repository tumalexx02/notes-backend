package login

import (
	"errors"
	"log/slog"
	resp "main/internal/http-server/api/response"
	"main/internal/models/user"
	"main/internal/storage"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type Response struct {
	resp.Response
	Token string `json:"token"`
}

type UserGetter interface {
	GetUser(email string) (user.User, error)
}

func New(log *slog.Logger, userGetter UserGetter, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.login.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to decode request body"))

			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request body", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("invalid request body"))

			return
		}

		user, err := userGetter.GetUser(req.Email)
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.Attr{Key: "email", Value: slog.StringValue(req.Email)})

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}
		if err != nil {
			log.Error("failed to get user", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to get user"))

			return
		}

		if !checkPassword(req.Password, user.PasswordHash) {
			log.Error("invalid password", slog.Attr{Key: "email", Value: slog.StringValue(req.Email)})

			render.JSON(w, r, resp.Error("invalid password"))

			return
		}

		// TODO: add configurable expiration time
		// TODO: add refresh token
		_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})
		if err != nil {
			log.Error("failed to encode token", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to encode token"))

			return
		}

		log.Info("success login", slog.Attr{Key: "email", Value: slog.StringValue(req.Email)})

		render.JSON(w, r, Response{resp.OK(), tokenString})
	}
}

func checkPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
