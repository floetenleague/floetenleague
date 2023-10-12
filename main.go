package main

import (
	"os"
	"time"

	"github.com/floetenleague/floetenleague/api"
	"github.com/floetenleague/floetenleague/config"
	"github.com/floetenleague/floetenleague/database"
	"github.com/floetenleague/floetenleague/dtest"
	"github.com/floetenleague/floetenleague/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.InfoLevel)
	log.Info().Msg("Logger initialized")
	cfg := config.Parse(os.Args[1])
	db := database.New(cfg)

	if cfg.DebugInsert {
		dtest.Test(db)
	}

	app := api.New(cfg, db)
	log.Info().Str("addr", cfg.ServerAddr).Msg("listening")
	err := server.Start(app, cfg.ServerAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("server")
	}
}
