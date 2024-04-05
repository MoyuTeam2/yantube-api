package main

import (
	"api/config"
	"api/server"
	"flag"

	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()

	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	log.Debug().Any("config", config.Config).Msg("config file loaded successfully")

	server.Start()

}
