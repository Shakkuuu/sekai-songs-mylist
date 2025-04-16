package config

import (
	"github.com/cockroachdb/errors"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBHost         string `env:"DB_HOST"`
	DBUserName     string `env:"DB_USER_NAME"`
	DBUserPassword string `env:"DB_USER_PASSWORD"`
	DBName         string `env:"DB_NAME"`
	DBPort         int    `env:"DB_PORT"`
	ServerPort     int    `env:"SERVER_PORT"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return cfg, nil
}
