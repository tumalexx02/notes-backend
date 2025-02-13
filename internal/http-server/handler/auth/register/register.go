package register

import (
	"errors"
	"log/slog"
	"main/internal/auth"
	"main/internal/config"
	resp "main/internal/http-server/api/response"
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
	Name     string `json:"name" validate:"required,max=31"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type Response struct {
	resp.Response
	auth.Tokens `json:"tokens"`
}

type Register interface {
	CreateUser(email, name, password string) (string, error)
	auth.RefreshTokenCreator
}

func New(cfg *config.Config, log *slog.Logger, register Register, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.register.New"

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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to hash password"))

			return
		}

		userID, err := register.CreateUser(req.Email, req.Name, string(hashedPassword))
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Error("user already exists", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}
		if err != nil {
			log.Error("failed to create user", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to create user"))

			return
		}

		refreshExp := time.Now().Add(cfg.Authorization.RefreshTTL)
		accessExp := time.Now().Add(cfg.Authorization.AccessTTL)

		tokens, err := auth.GenerateTokens(register, userID, refreshExp, accessExp, tokenAuth)
		if err != nil {
			log.Error("failed to generate tokens", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to generate tokens"))

			return
		}

		log.Info("user registered", slog.String("user_id", userID))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Tokens:   tokens,
		})
	}
}
