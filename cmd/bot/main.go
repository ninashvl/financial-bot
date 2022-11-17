package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"golang.org/x/sys/unix"

	"gitlab.ozon.dev/ninashvl/homework-1/config"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/clients/tg"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/messages"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		unix.SIGTERM, unix.SIGKILL, unix.SIGINT)
	defer cancel()
	logger := zerolog.New(os.Stdout)

	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("config init failed")
	}

	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("tg client init failed")
	}

	bot := messages.New(tgClient, cfg, logger)
	tgClient.ListenUpdates(ctx, bot)
	logger.Info().Msg("application gracefully stopped")
}
