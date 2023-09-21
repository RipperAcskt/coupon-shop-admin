package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDBUsername string `env:"POSTGRES_DB_USERNAME"`
	PostgresDBPassword string `env:"POSTGRES_DB_PASSWORD"`
	PostgresDBHost     string `env:"POSTGRES_DB_HOST"`
	PostgresDBName     string `env:"POSTGRES_DB_NAME"`
	MigratePath        string `env:"MIGRATE_PATH"`

	ServerHost string `env:"SERVER_HOST"`
}

func New() (Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		return Config{}, fmt.Errorf("load failed: %w", err)
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse failed: %w", err)
	}
	return cfg, nil
}

func (c *Config) GetPostgresUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s", c.PostgresDBUsername, c.PostgresDBPassword, c.PostgresDBHost, c.PostgresDBName)
}
