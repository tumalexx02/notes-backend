package router

import (
	"log/slog"
	"main/internal/config"
	"main/internal/http-server/handler/auth/login"
	"main/internal/http-server/handler/auth/me"
	"main/internal/http-server/handler/auth/refresh"
	"main/internal/http-server/handler/auth/register"
	"main/internal/http-server/middleware/authenticator"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

type Authorizer interface {
	register.Register
	login.Loginer
	refresh.RefreshTokener
}

func (r *Router) InitAuthRoutes(storage Storage, logger *slog.Logger, cfg *config.Config) {
	r.Route("/user", func(userRouter chi.Router) {
		userRouter.Post("/register", register.New(cfg, logger, storage, r.jwtauth))
		userRouter.Post("/login", login.New(cfg, logger, storage, r.jwtauth))
		userRouter.Post("/refresh", refresh.New(cfg, logger, storage, r.jwtauth))

		userRouter.Group(func(protected chi.Router) {
			protected.Use(jwtauth.Verifier(r.jwtauth))
			protected.Use(authenticator.Authenticator(r.jwtauth, logger))

			protected.Get("/me", me.New(logger, r.jwtauth))
		})
	})
}

func generateAuthToken(cfg *config.Config) *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return tokenAuth
}
