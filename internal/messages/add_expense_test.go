package messages

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "gitlab.ozon.dev/ninashvl/homework-1/internal/mocks/messages"
	statestorage "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/dialogue_state_storage/mocks"
	expstorage "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage/mocks"
)

func TestBot_addExpenseInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	expStore := expstorage.NewMockIStorage(ctrl)
	stStore := statestorage.NewMockIStorage(ctrl)

	bot := &Bot{
		tgClient:        sender,
		expStorage:      expStore,
		dlgStateStorage: stStore,
	}
	userID := int64(1)

	// case 1 - невалидный формат сообщения для парсинга
	msg := &Message{
		Text:   "1234",
		UserID: userID,
	}
	sender.EXPECT().SendMessage(gomock.Any(), invalidMsg, userID)

	err := bot.AddExpense(context.TODO(), msg)
	assert.Nil(t, err, "addExpense error")

	// case 2 - невалидное число для парсинга
	msg = &Message{
		Text:   "сумма,категория",
		UserID: userID,
	}
	sender.EXPECT().SendMessage(gomock.Any(), invalidMsg, userID)

	err = bot.AddExpense(context.TODO(), msg)
	assert.Nil(t, err, "addExpense error")

	// case 2 - невалидное число для парсинга
	msg = &Message{
		Text:   "1,категория, 10.11.2001",
		UserID: userID,
	}
	sender.EXPECT().SendMessage(gomock.Any(), invalidTimestamp, userID)

	err = bot.AddExpense(context.TODO(), msg)
	assert.Nil(t, err, "addExpense error")
}

func TestBot_addExpenseSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	expStore := expstorage.NewMockIStorage(ctrl)
	stStore := statestorage.NewMockIStorage(ctrl)
	bot := &Bot{
		tgClient:        sender,
		expStorage:      expStore,
		dlgStateStorage: stStore,
	}
	userID := int64(1)
	expStore.EXPECT().Add(gomock.Any(), userID, gomock.Any())
	sender.EXPECT().SendMessage(gomock.Any(), savedMsg, userID)
	expStore.EXPECT().GetLimit(gomock.Any(), userID)
	msg := &Message{
		Text:   "1,категория",
		UserID: userID,
	}
	err := bot.AddExpense(context.TODO(), msg)
	assert.Nil(t, err, "addExpense error")
}
