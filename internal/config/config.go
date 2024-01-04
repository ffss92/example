package config

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/caarlos0/env/v10"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

type Config struct {
	Port         int    `env:"APP_PORT" envDefault:"4000"`
	Env          string `env:"APP_ENV" envDefault:"development"`
	SecretKey    Secret `env:"SECRET_KEY,required"`
	DatabasePath string `env:"DATABASE_PATH" envDefault:"example.db"`
}

func New() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Writes the current config to the specificied writer as JSON.
func (c Config) Dump(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}

// Checks if the current config is in development mode.
func (c Config) IsDev() bool {
	return c.Env == EnvDevelopment
}

// Checks if the current config is in production mode.
func (c Config) IsProd() bool {
	return c.Env == EnvProduction
}

// Returns the address to listen on.
func (c Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}
