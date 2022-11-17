package messages

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
)

func (s *Bot) AddExpense(ctx context.Context, msg *Message) error {
	s.logger.Debug().Str("text", msg.Text).Int64("user", msg.UserID).Msg("AddExpense func started")
	parts := strings.Split(msg.Text, ",")
	if len(parts) < 2 {
		s.logger.Error().Str("text", msg.Text).Int64("user", msg.UserID).Msg("Len parts more than 2")
		return s.tgClient.SendMessage(invalidMsg, msg.UserID)
	}
	num, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		s.logger.Error().Str("text", msg.Text).Int64("user", msg.UserID).Err(err)
		return s.tgClient.SendMessage(invalidMsg, msg.UserID)
	}
	category := strings.TrimSpace(parts[1])
	exp := &models.Expense{
		Amount:   num,
		Category: category,
		Date:     time.Now(),
	}
	// check on existing timestamp
	if len(parts) > 2 {
		t, err := time.Parse("2006-01-02", strings.TrimSpace(parts[2]))
		if err != nil {
			s.logger.Error().Str("text", msg.Text).Int64("user", msg.UserID).Err(err).Msg("Parsefloat of msg error")
			return s.tgClient.SendMessage(invalidTimestamp, msg.UserID)
		}
		exp.Date = t
	}

	err = s.expStorage.Add(ctx, msg.UserID, exp)
	if err != nil {
		s.logger.Error().Str("text", msg.Text).Int64("user", msg.UserID).Err(err)
		return err
	}

	_ = s.checkLimit(ctx, msg.UserID)

	return s.tgClient.SendMessage(savedMsg, msg.UserID)
}

func (s *Bot) checkLimit(ctx context.Context, userID int64) error {
	limit, err := s.expStorage.GetLimit(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) || limit == 0 {
		return nil
	}

	res, err := s.expStorage.GetByRange(ctx, userID, expense_storage.Month)
	if err != nil {
		s.logger.Error().Err(err)
		return s.tgClient.SendMessage(err.Error(), userID)
	}

	count := 0.
	for i := 0; i < len(res); i++ {
		count += res[i].Amount
	}

	curr, err := s.expStorage.GetCurrency(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Send()
		return s.tgClient.SendMessage(err.Error(), userID)
	}

	if count > limit {
		msg := fmt.Sprintf("У вас превышен лимит!\nлимит - %s %s\nтраты за месяц - %s %s",
			strconv.FormatFloat(limit, 'f', 2, 64), curr,
			strconv.FormatFloat(count, 'f', 2, 64), curr)
		return s.tgClient.SendMessage(msg, userID)
	}

	return nil
}
