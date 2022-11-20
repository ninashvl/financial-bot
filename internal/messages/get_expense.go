package messages

import (
	"context"
	"strconv"
	"strings"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *Bot) GetExpense(ctx context.Context, msg *Message) error {
	var span trace.Span
	ctx, span = otel.Tracer("update").Start(ctx, "bot.GetExpense")
	defer span.End()

	s.logger.Debug().Str("text", msg.Text).Int64("user", msg.UserID).Msg("GetExpense func started")
	var res []*models.TotalExpense
	var err error
	switch strings.TrimSpace(msg.Text) {
	case "День":
		res, err = s.expStorage.GetByRange(ctx, msg.UserID, expense_storage.Day)
	case "Месяц":
		res, err = s.expStorage.GetByRange(ctx, msg.UserID, expense_storage.Month)
	case "Год":
		res, err = s.expStorage.GetByRange(ctx, msg.UserID, expense_storage.Year)
	default:
		return s.tgClient.SendMessage(ctx, invalidRange, msg.UserID)
	}
	if err != nil {
		return s.tgClient.SendMessage(ctx, err.Error(), msg.UserID)
	}
	if len(res) == 0 {
		span.SetStatus(codes.Error, "empty res")
		return s.tgClient.SendMessage(ctx, expensesNotFound, msg.UserID)
	}
	curr, err := s.expStorage.GetCurrency(ctx, msg.UserID)
	if err != nil {
		return err
	}
	builder := strings.Builder{}
	for _, v := range res {
		builder.WriteString(v.Category)
		builder.WriteString(": ")
		builder.WriteString(strconv.FormatFloat(v.Amount, 'f', 2, 64))
		builder.WriteString(" " + curr)
		builder.WriteString("\n")
	}
	s.logger.Info().Str("text", msg.Text).Int64("user", msg.UserID).Msg("GetExpense func executed")
	return s.tgClient.SendMessage(ctx, builder.String(), msg.UserID)
}
