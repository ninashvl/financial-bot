package messages

import (
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
)

func (s *Bot) SetCurrency(msg *Message) error {
	if msg.Text == models.UsdCurrency || msg.Text == models.RubCurrency ||
		msg.Text == models.CnyCurrency || msg.Text == models.EurCurrency {
		s.expStorage.SetCurrency(msg.UserID, msg.Text)
		return s.tgClient.SendMessage("Валюта трат сохранена", msg.UserID)
	}
	return s.tgClient.SendMessage("Невалидная валюта", msg.UserID)
}
