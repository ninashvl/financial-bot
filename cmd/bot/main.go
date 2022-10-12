package main

import (
	"log"

	"gitlab.ozon.dev/ninashvl/homework-1/config"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/clients/tg"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/messages"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(config)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	msgModel := messages.New(tgClient)

	tgClient.ListenUpdates(msgModel)
}
