package main

import (
	"fmt"
	"main/internal/storage/postgres"
)

func main() {
	fmt.Println("Hello, go!")

	storage, err := postgres.New()
	if err != nil {
		panic(err)
	}

	_ = storage
}
