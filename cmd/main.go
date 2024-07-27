package main

import (
	"chat-backend/internal/app"
	"chat-backend/internal/config"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.ReadConfig()
	ctx := context.Background()

	app := app.New(cfg)
	go func() {
		if err := app.Start(ctx); err != nil {
			panic(err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-exit
	fmt.Printf("Shutting down...\n")

	app.Shutdown(ctx)
}
