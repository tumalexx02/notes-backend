package main

import (
	"log/slog"
	"main/internal/app"
	"main/internal/config"
	"main/internal/logger"
	"os"
)

// TODO: change revoke logic to delete refresh token
// TODO: custom error messages (for separate feature)

func main() {
	// config init
	cfg := config.MustLoad()

	// logger init``
	log := logger.New(cfg)

	// app init
	myApp, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to create app", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}

	// app start
	err = myApp.Start()
	if err != nil {
		log.Error("failed to start app", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}
}
