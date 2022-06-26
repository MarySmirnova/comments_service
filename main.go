package main

import (
	"time"

	"github.com/MarySmirnova/comments_service/internal"
	"github.com/MarySmirnova/comments_service/internal/config"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

var cfg config.Application

func init() {
	gotenv.Load(".env")
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(lvl)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp,
	})

}

func main() {
	app, err := internal.NewApplication(cfg)
	if err != nil {
		panic(err)
	}

	app.StartServer()
}
