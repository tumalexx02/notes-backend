package login

import (
	"errors"
	"log/slog"
	"main/internal/auth"
	"main/internal/config"
	resp "main/internal/http-server/api/response"
	"main/internal/models/user"
	"main/internal/storage"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
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
	Tokens auth.Tokens `json:"tokens"`
}

type Loginer interface {
	GetUser(email string) (user.User, error)
	auth.RefreshTokenCreator
}

func New(cfg *config.Config, log *slog.Logger, loginer Loginer, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
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

		userFromDb, err := loginer.GetUser(req.Email)
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

		if !checkPassword(req.Password, userFromDb.PasswordHash) {
			log.Error("invalid password", slog.Attr{Key: "email", Value: slog.StringValue(req.Email)})

			render.JSON(w, r, resp.Error("invalid password"))

			return
		}

		refreshExp := time.Now().Add(cfg.Authorization.RefreshTTL)
		accessExp := time.Now().Add(cfg.Authorization.AccessTTL)

		tokens, err := auth.GenerateTokens(loginer, userFromDb.ID, refreshExp, accessExp, tokenAuth)
		if err != nil {
			log.Error("failed to generate tokens", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to generate tokens"))

			return
		}

		log.Info("success login", slog.Attr{Key: "email", Value: slog.StringValue(req.Email)})

		render.JSON(w, r, Response{resp.OK(), tokens})
	}
}

func checkPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
