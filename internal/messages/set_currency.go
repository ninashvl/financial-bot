package messages

import (
	"context"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func (s *Bot) SetCurrency(ctx context.Context, msg *Message) error {
	var span trace.Span
	ctx, span = otel.Tracer("update").Start(ctx, "bot.SetCurrency")
	defer span.End()

	s.logger.Debug().Str("text", msg.Text).Int64("user", msg.UserID).Msg("Set currency func started")
	if msg.Text == models.UsdCurrency || msg.Text == models.RubCurrency ||
		msg.Text == models.CnyCurrency || msg.Text == models.EurCurrency {
		err := s.expStorage.SetCurrency(ctx, msg.UserID, msg.Text)
		if err != nil {
			s.logger.Error().Str("text", msg.Text).Int64("user", msg.UserID).Err(err).Msg("SetCurrency error")
			return err
		}
		s.logger.Debug().Str("text", msg.Text).Int64("user", msg.UserID).Msg("Set currency func executed")
		return s.tgClient.SendMessage(ctx, currencySaved, msg.UserID)
	}
	return s.tgClient.SendMessage(ctx, invalidCurrency, msg.UserID)
}
