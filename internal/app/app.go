package app

import (
	"chat-backend/internal/config"
	"chat-backend/internal/handler/http"
	"chat-backend/internal/jwt"
	"chat-backend/internal/logger"
	"chat-backend/internal/repository"
	"chat-backend/internal/repository/postgresql"
	"chat-backend/internal/service"
	"chat-backend/pkg/clock"
	"context"
)

type App struct {
	services   *service.Services
	repository *repository.Repository
	handler    *http.Handler
	logger     logger.Logger
}

func New(config config.Config) *App {
	logger := logger.NewLogrusLogger(config.App.LogLevel, config.App.IsDev)
	logger.Infof("config loaded")

	if config.App.IsTesting {
		clock.InitClock(true)
	}

	logger.Infof("connecting to postgresql on %s:%d", config.Postgresql.Host, config.Postgresql.Port)
	psql, err := postgresql.New(config.Postgresql)
	if err != nil {
		logger.Fatalf("connect to postgresql: %s", err)
	}

	err = postgresql.Migrate(config.Postgresql)
	if err != nil {
		logger.Warnf("migrate database: %s", err)
	}

	jwt := jwt.New(config.JWT)
	repository := repository.New(psql, logger)
	services := service.New(repository, jwt)
	handler := http.New(config.Server, services, logger, jwt)

	return &App{
		services:   services,
		repository: repository,
		handler:    handler,
		logger:     logger,
	}
}

func (app *App) Start() error {
	if err := app.handler.Start(); err != nil {
		app.logger.Errorf("failed to start app: %s", err)
		return err
	}

	return nil
}

func (app *App) Shutdown(ctx context.Context) {
	if err := app.handler.Stop(ctx); err != nil {
		panic(err)
	}
}
