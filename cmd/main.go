package main

import (
	"main/internal/app"
	"main/internal/config"
	"main/internal/logger"
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	cfg := config.MustLoad()

	log := logger.InitLogger(cfg)

	myApp, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to create app", slog.Attr{"error", slog.StringValue(err.Error())})
		os.Exit(1)
	}

	_ = myApp
}
