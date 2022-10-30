package messages

import (
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/dialogue_state_storage"
	in_mem_dlg "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/dialogue_state_storage/in_memory"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
	in_mem_exp "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage/in_memory"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
	SendRangeKeyboard(userID int64, text string) error
}

type Bot struct {
	tgClient        MessageSender
	expStorage      expense_storage.IStorage
	dlgStateStorage dialogue_state_storage.IStorage
}

func New(tgClient MessageSender) *Bot {
	return &Bot{
		tgClient:        tgClient,
		expStorage:      in_mem_exp.New(),
		dlgStateStorage: in_mem_dlg.New(),
	}
}

type Message struct {
	Text      string
	UserID    int64
	IsCommand bool
}

func (s *Bot) HandleCommand(msg *Message) error {
	switch {
	case msg.Text == startCommand:
		return s.tgClient.SendMessage("hello", msg.UserID)
	case msg.Text == helpCommand && msg.IsCommand:
		return s.tgClient.SendMessage(help, msg.UserID)
	case msg.Text == addCommand && msg.IsCommand:
		s.dlgStateStorage.Add(msg.UserID, models.AddCommandState)
		return s.tgClient.SendMessage(addMessage, msg.UserID)
	case msg.Text == getExpensesCommand && msg.IsCommand:
		s.dlgStateStorage.Add(msg.UserID, models.GetCommandState)
		return s.tgClient.SendRangeKeyboard(msg.UserID, "Выберите диапазон")
	default:
		return s.tgClient.SendMessage("не знаю эту команду", msg.UserID)
	}
}

func (s *Bot) HandleMessage(msg *Message) error {
	switch {
	case !msg.IsCommand && s.dlgStateStorage.Get(msg.UserID) == models.AddCommandState:
		return s.addExpense(msg)
	case !msg.IsCommand && s.dlgStateStorage.Get(msg.UserID) == models.GetCommandState:
		return s.GetExpense(msg)
	default:
		return s.tgClient.SendMessage("воспользуйтесь /help", msg.UserID)
	}
}

func (s *Bot) IncomingMessage(msg *Message) error {
	defer func() {
		if !msg.IsCommand {
			s.dlgStateStorage.DeleteState(msg.UserID)
		}
	}()
	if msg.IsCommand {
		return s.HandleCommand(msg)
	}
	return s.HandleMessage(msg)
}
