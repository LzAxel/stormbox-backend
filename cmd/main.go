package main

import (
	"chat-backend/internal/app"
	"chat-backend/internal/config"
)

func main() {
	cfg := config.ReadConfig()

	app := app.New(cfg)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
