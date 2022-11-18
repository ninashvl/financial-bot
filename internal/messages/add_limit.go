package messages

import (
	"context"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func (s *Bot) AddLimit(ctx context.Context, msg *Message) error {
	var span trace.Span
	ctx, span = otel.Tracer("update").Start(ctx, "bot.AddLimit")
	defer span.End()

	limit, err := strconv.ParseFloat(msg.Text, 64)
	if err != nil {
		s.logger.Error().Err(err)
		return s.tgClient.SendMessage(ctx, invalidMsg, msg.UserID)
	}

	err = s.expStorage.SetLimit(ctx, msg.UserID, limit)
	if err != nil {
		s.logger.Error().Err(err)
		return s.tgClient.SendMessage(ctx, err.Error(), msg.UserID)
	}
	return s.tgClient.SendMessage(ctx, limitSuccessfulSet, msg.UserID)
}
