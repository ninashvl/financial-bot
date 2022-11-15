package messages

import (
	"context"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
)

func (s *Bot) SetCurrency(ctx context.Context, msg *Message) error {
	if msg.Text == models.UsdCurrency || msg.Text == models.RubCurrency ||
		msg.Text == models.CnyCurrency || msg.Text == models.EurCurrency {
		err := s.expStorage.SetCurrency(ctx, msg.UserID, msg.Text)
		if err != nil {
			return err
		}
		return s.tgClient.SendMessage(currencySaved, msg.UserID)
	}
	return s.tgClient.SendMessage(invalidCurrency, msg.UserID)
}
