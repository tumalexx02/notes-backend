package register

import (
	"errors"
	"log/slog"
	"main/internal/auth"
	"main/internal/config"
	resp "main/internal/http-server/api/response"
	resperrors "main/internal/http-server/api/response-errors"
	"main/internal/http-server/api/validate"
	"main/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
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

		if err := validate.DecodeAndValidateRequestJson(&req, w, r, log); err != nil {
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError))

			return
		}

		userID, err := register.CreateUser(req.Email, req.Name, string(hashedPassword))
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Error("user already exists", "error", err)

			render.JSON(w, r, resp.Error(resperrors.ErrUserIsAlreadyExists))

			return
		}
		if err != nil {
			log.Error("failed to create user", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError))

			return
		}

		tokens, err := auth.GenerateTokens(userID, register, cfg, tokenAuth)
		if err != nil {
			log.Error("failed to generate tokens", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(resperrors.ErrInternalServerError))

			return
		}

		log.Info("user registered", slog.String("user_id", userID))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Tokens:   tokens,
		})
	}
}
