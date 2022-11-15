package messages

import (
	"context"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "gitlab.ozon.dev/ninashvl/homework-1/internal/mocks/messages"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	statestorage "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/dialogue_state_storage/mocks"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
	expstorage "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage/mocks"
)

func TestBot_getExpenseInvalid(t *testing.T) {
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

	// case 1 - invalid range
	sender.EXPECT().SendMessage(invalidRange, userID)
	msg := &Message{
		Text:   "",
		UserID: userID,
	}
	err := bot.GetExpense(context.TODO(), msg)
	assert.Nil(t, err, "GetExpense Error")

	// case 2 - empty result
	expStore.EXPECT().GetByRange(context.TODO(), userID, expense_storage.Day).Return([]*models.TotalExpense{}, nil)
	msg = &Message{
		Text:   "День",
		UserID: userID,
	}
	sender.EXPECT().SendMessage(expensesNotFound, userID)
	err = bot.GetExpense(context.TODO(), msg)
	assert.Nil(t, err, expensesNotFound)
}

func TestBot_getExpense(t *testing.T) {
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
	amount := 100.
	category := "test"
	curr := "RUB"

	expStore.EXPECT().GetByRange(context.TODO(), userID, expense_storage.Day).Return([]*models.TotalExpense{
		{amount, category},
	}, nil)
	expStore.EXPECT().GetCurrency(context.TODO(), userID).Return("RUB", nil)
	resMsg := category + ": " + strconv.FormatFloat(amount, 'f', 2, 64) + " " + curr + "\n"
	sender.EXPECT().SendMessage(resMsg, userID)

	msg := &Message{
		Text:   "День",
		UserID: userID,
	}
	err := bot.GetExpense(context.TODO(), msg)
	assert.Nil(t, err, expensesNotFound)
}
