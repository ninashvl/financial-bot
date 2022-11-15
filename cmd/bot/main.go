package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"

	"golang.org/x/sys/unix"

	"gitlab.ozon.dev/ninashvl/homework-1/config"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/clients/tg"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/messages"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		unix.SIGTERM, unix.SIGKILL, unix.SIGINT)
	defer cancel()
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(cfg)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	bot := messages.New(tgClient, cfg)
	tgClient.ListenUpdates(ctx, bot)
	fmt.Println("[INFO] application gracefully stopped")
}
