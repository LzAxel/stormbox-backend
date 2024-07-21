package config

import (
	"sync"

	"chat-backend/internal/handler/http"
	"chat-backend/internal/jwt"
	"chat-backend/internal/repository/postgresql"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "configs/dev.yaml"
)

type AppConfig struct {
	IsDev     bool   `yaml:"isDev" env:"IS_DEV"`
	IsTesting bool   `yaml:"isTesting" env:"IS_TESTING"`
	LogLevel  string `yaml:"logLevel" env:"LOG_LEVEL"`
}

type Config struct {
	Postgresql postgresql.Config `yaml:"postgres"`
	Server     http.Config       `yaml:"server"`
	App        AppConfig         `yaml:"app"`
	JWT        jwt.Config        `yaml:"jwt"`
}

var (
	config Config
	once   sync.Once
)

func ReadConfig() Config {
	once.Do(func() {
		if err := cleanenv.ReadConfig(configPath, &config); err != nil {
			panic(err)
		}
	})

	return config
}
