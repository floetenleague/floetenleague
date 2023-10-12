package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DBConnection    string `env:"FL_DB_CONNECTION" envDefault:"postgres://root:root@localhost/floetenleague?sslmode=disable"`
	POEClientID     string `env:"FL_POE_CLIENT_ID"`
	POEClientSecret string `env:"FL_POE_CLIENT_SECRET"`
	SessionKey      string `env:"FL_SESSION_KEY"`
	DebugInsert     bool   `env:"FL_DEBUG_INSERT"`
	ServerAddr      string `env:"FL_SERVER_ADDR" envDefault:"127.0.0.1:8080"`
}

func Parse(path string) *Config {
	_ = godotenv.Load(path)
	var c Config
	err := env.Parse(&c)
	if err != nil {
		log.Fatal().Err(err).Msg("parse config")
	}

	return &c
}
