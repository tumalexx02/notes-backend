package logger

import (
	"log/slog"
	"main/internal/config"
	"os"
)

// init logger by env field from config
func New(cfg *config.Config) *slog.Logger {
	var log *slog.Logger

	switch cfg.Env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
