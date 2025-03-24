package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Mode     string `env:"MODE"`
	Postgres Postgres
	HTTP     HTTP
}

type HTTP struct {
	Port         int           `env:"HTTP_PORT" default:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" default:"10s"`
}

type Postgres struct {
	DSN      string `env:"POSTGRES_DSN" default:"postgresql://postgres:dev@localhost:5432/postgres?sslmode=disable"`
	MaxConns int32  `env:"POSTGRES_MAX_CONNS" default:"30"`
}

func New(envPath string) (*Config, error) {
	var c Config

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("failed to load env file: %w", err)
	}

	if err := env.ParseWithOptions(&c, env.Options{RequiredIfNoDef: true}); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return &c, nil
}
