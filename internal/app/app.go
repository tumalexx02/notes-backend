package app

import (
	"log/slog"
	"main/internal/config"
	"main/internal/storage/postgres"
)

type App struct {
	Config  *config.Config
	Logger  *slog.Logger
	Storage *postgres.Storage
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	storage, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:  cfg,
		Logger:  log,
		Storage: storage,
	}, nil
}
