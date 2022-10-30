package messages

import (
	"strconv"
	"strings"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
)

func (s *Bot) GetExpense(msg *Message) error {
	var res []*models.TotalExpense
	var err error
	switch strings.TrimSpace(msg.Text) {
	case "День":
		res, err = s.expStorage.GetByRange(msg.UserID, expense_storage.Day)
	case "Месяц":
		res, err = s.expStorage.GetByRange(msg.UserID, expense_storage.Month)
	case "Год":
		res, err = s.expStorage.GetByRange(msg.UserID, expense_storage.Year)
	default:
		return s.tgClient.SendMessage(invalidRange, msg.UserID)
	}
	if err != nil {
		return s.tgClient.SendMessage(err.Error(), msg.UserID)
	}
	if len(res) == 0 {
		return s.tgClient.SendMessage(expensesNotFound, msg.UserID)
	}
	curr := s.expStorage.GetCurrency(msg.UserID)
	builder := strings.Builder{}
	for _, v := range res {
		builder.WriteString(v.Category)
		builder.WriteString(": ")
		builder.WriteString(strconv.FormatFloat(v.Amount, 'f', 2, 64))
		builder.WriteString(" " + curr)
		builder.WriteString("\n")
	}
	return s.tgClient.SendMessage(builder.String(), msg.UserID)
}
