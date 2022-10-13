package messages

import (
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/InMemory"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Bot struct {
	tgClient MessageSender
	storage  storage.IStorage
}

func New(tgClient MessageSender) *Bot {
	return &Bot{
		tgClient: tgClient,
		storage:  InMemory.New(),
	}
}

type Message struct {
	Text   string
	UserID int64
}

func addExpense(s *Bot, msg Message) error {
	parts := strings.Split(msg.Text, ",")
	if parts[0] != addCommand {
		return s.tgClient.SendMessage(invalidCommand, msg.UserID)
	}
	num, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return s.tgClient.SendMessage(invalidCommand+err.Error(), msg.UserID)
	}
	category := strings.TrimSpace(parts[2])
	s.storage.Add(msg.UserID, &models.Expense{
		Amount:   num,
		Category: category,
		Date:     time.Now(),
	})
	return s.tgClient.SendMessage("Сохранено", msg.UserID)
}

func (s *Bot) IncomingMessage(msg Message) error {
	switch {
	case msg.Text == "/start":
		return s.tgClient.SendMessage("hello", msg.UserID)
	case msg.Text == "/help":
		return s.tgClient.SendMessage(help, msg.UserID)
	case strings.Contains(msg.Text, addCommand):
		// Формат ввода команды пользователем: 📎, сумма, категория
		return addExpense(s, msg)
	case strings.Contains(msg.Text, getExpensesCommand):
		// 📤, num
		parts := strings.Split(msg.Text, ",")
		if len(parts) != 2 {
			return s.tgClient.SendMessage(invalidCommand, msg.UserID)
		}
		if parts[0] != getExpensesCommand {
			return s.tgClient.SendMessage(invalidCommand, msg.UserID)
		}
		var res []*models.TotalExpense
		var err error
		switch strings.TrimSpace(parts[1]) {
		case "1":
			res, err = s.storage.GetByRange(msg.UserID, storage.Day)
		case "2":
			res, err = s.storage.GetByRange(msg.UserID, storage.Month)
		case "3":
			res, err = s.storage.GetByRange(msg.UserID, storage.Year)
		default:
			return s.tgClient.SendMessage(invalidRange, msg.UserID)
		}
		if err != nil {
			return s.tgClient.SendMessage(err.Error(), msg.UserID)
		}
		if len(res) == 0 {
			return s.tgClient.SendMessage(expensesNotFound, msg.UserID)
		}
		builder := strings.Builder{}
		for _, v := range res {
			builder.WriteString(v.Category)
			builder.WriteString(": ")
			builder.WriteString(strconv.FormatFloat(v.Amount, 'f', -2, 64))
			builder.WriteString("\n")
		}
		return s.tgClient.SendMessage(builder.String(), msg.UserID)
	}

	return s.tgClient.SendMessage("не знаю эту команду", msg.UserID)
}
