package main

import (
	"log/slog"
	"main/internal/app"
	"main/internal/config"
	"main/internal/logger"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg)

	myApp, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to create app", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}

	err = myApp.Start()
	if err != nil {
		log.Error("failed to start app", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}
}
