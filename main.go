package main

import (
	"api/config"
	"api/server"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()

	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	log.Debug().Any("config", config.Config).Msg("config file loaded successfully")

	go server.StartHttp()
	go server.StartGrpc()

	// listen signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info().Msg("server shutdown")

}
