package main

import (
	"main/internal/app"
	"main/internal/config"
	"main/internal/logger"
	"os"
)

func main() {
	// config init
	cfg := config.MustLoad()

	// logger init
	log := logger.New(cfg)

	// app init
	myApp, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to create app", "error", err)
		os.Exit(1)
	}

	// app start
	err = myApp.Start()
	if err != nil {
		log.Error("failed to start app", "error", err)
		os.Exit(1)
	}
}
