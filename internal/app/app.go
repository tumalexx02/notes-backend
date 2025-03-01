package app

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/config"
	"main/internal/router"
	"main/internal/storage/postgres"
	"net/http"
)

type App struct {
	config  *config.Config
	logger  *slog.Logger
	storage *postgres.Storage
	router  *router.Router
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	const op = "app.Start"

	// init storage
	storage, err := postgres.New(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// init router and routes
	router := router.New(cfg, log)
	router.InitRoutes(storage, log, cfg)

	return &App{
		config:  cfg,
		logger:  log,
		storage: storage,
		router:  router,
	}, nil
}

func (a *App) Start() error {
	const op = "app.Start"

	// init context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init jobs
	startTokensRevokingJob(ctx, a.logger, a.storage)

	a.logger.Info("starting server", slog.String("address", a.config.HTTPServer.Address))

	// init server
	srv := &http.Server{
		Addr:         a.config.HTTPServer.Address,
		Handler:      a.router,
		ReadTimeout:  a.config.HTTPServer.Timeout,
		IdleTimeout:  a.config.HTTPServer.IdleTimeout,
		WriteTimeout: a.config.HTTPServer.Timeout,
	}

	// start server
	return fmt.Errorf("%s: %w", op, srv.ListenAndServe())
}
