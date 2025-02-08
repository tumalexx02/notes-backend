package main

import (
	"fmt"
	"main/internal/config"
	"main/internal/storage/postgres"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)

	storage, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}

	_ = storage
}
