package main

import (
	"api/config"
	"api/db"
	"api/server"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "api/pkg/memdump"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()

	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	zerolog.TimeFieldFormat = time.DateTime
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	consoleWriter.TimeFormat = time.DateTime
	file, err := os.OpenFile("apiserver.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, file)
	logger := zerolog.New(multi).With().Timestamp().Caller().Logger()
	log.Logger = logger

	log.Debug().Any("config", config.Config).Msg("config file loaded successfully")

	err = db.Init(config.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init db")
	}

	go server.StartHttp()
	go server.StartGrpc()

	// listen signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info().Msg("server shutdown")
}
