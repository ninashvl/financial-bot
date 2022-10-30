package messages

import (
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
)

func (s *Bot) addExpense(msg *Message) error {
	parts := strings.Split(msg.Text, ",")
	if len(parts) < 2 {
		return s.tgClient.SendMessage(invalidMsg, msg.UserID)
	}
	num, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return s.tgClient.SendMessage(invalidMsg+err.Error(), msg.UserID)
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
			return s.tgClient.SendMessage(invalidTimestamp+err.Error(), msg.UserID)
		}
		exp.Date = t
	}

	s.expStorage.Add(msg.UserID, exp)
	return s.tgClient.SendMessage("Сохранено", msg.UserID)
}
