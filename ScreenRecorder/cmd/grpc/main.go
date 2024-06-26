package main

import (
	"myapp/config"
	"myapp/internal/grpc/app"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	// Настройка логгера
	output := zerolog.ConsoleWriter{
		TimeFormat: "02.01.2006 15:04:05",
		Out:        os.Stdout,
	}
	log.Logger = log.Output(output)

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Msgf("Config error: %s", err)
	}

	app.Run(cfg)
}
