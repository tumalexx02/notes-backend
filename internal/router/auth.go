package router

import (
	"log/slog"
	"main/internal/config"
	"main/internal/http-server/handler/auth/login"
	"main/internal/http-server/handler/auth/me"
	"main/internal/http-server/handler/auth/register"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

func (r *Router) InitAuthRoutes(storage Storage, logger *slog.Logger, cfg *config.Config) {
	r.Route("/user", func(userRouter chi.Router) {
		userRouter.Post("/register", register.New(logger, storage, r.jwtauth))
		userRouter.Post("/login", login.New(logger, storage, r.jwtauth))

		userRouter.Group(func(protected chi.Router) {
			protected.Use(jwtauth.Verifier(r.jwtauth))
			protected.Use(jwtauth.Authenticator(r.jwtauth))

			protected.Get("/me", me.New(logger, r.jwtauth))
		})
	})
}

func generateAuthToken(cfg *config.Config) *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return tokenAuth
}
