package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/metrics"
	"go.opentelemetry.io/otel"
	"golang.org/x/sys/unix"

	"gitlab.ozon.dev/ninashvl/homework-1/config"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/clients/tg"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/messages"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
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

	trace, err := tracerProvider(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("trace provider fatal error")
	}
	otel.SetTracerProvider(trace)
	tgClient, err := tg.New(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("tg client init failed")
	}
	go metrics.StartMetricServer()
	bot := messages.New(tgClient, cfg, logger)
	tgClient.ListenUpdates(ctx, bot)
	logger.Info().Msg("application gracefully stopped")
}

func tracerProvider(cfg *config.Service) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.TraceUrl())))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fin_bot"),
			attribute.String("environment", "dev"),
		)),
	)
	return tp, nil
}
