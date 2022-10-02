package main

import (
	"log"

	"gitlab.ozon.dev/lukyantsev-pa/lukyantsev-pavel/internal/clients/tg"
	"gitlab.ozon.dev/lukyantsev-pa/lukyantsev-pavel/internal/config"
	"gitlab.ozon.dev/lukyantsev-pa/lukyantsev-pavel/internal/model/messages"
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
