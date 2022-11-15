package messages

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "gitlab.ozon.dev/ninashvl/homework-1/internal/mocks/messages"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	statestorage "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/dialogue_state_storage/mocks"
	expstorage "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage/mocks"
)

func TestBot_SetCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	expStore := expstorage.NewMockIStorage(ctrl)
	stStore := statestorage.NewMockIStorage(ctrl)

	bot := &Bot{
		tgClient:        sender,
		expStorage:      expStore,
		dlgStateStorage: stStore,
	}

	// case 1 - invalid income currency
	userID := int64(1)
	curr := "ARS"
	sender.EXPECT().SendMessage(invalidCurrency, userID)
	msg := &Message{
		Text:   curr,
		UserID: userID,
	}

	err := bot.SetCurrency(context.TODO(), msg)
	assert.Nil(t, err, "SetCurrency error")

	// case 2 - success
	curr = models.UsdCurrency
	expStore.EXPECT().SetCurrency(context.TODO(), userID, curr)
	sender.EXPECT().SendMessage(currencySaved, userID)

	msg = &Message{
		Text:   curr,
		UserID: userID,
	}
	err = bot.SetCurrency(context.TODO(), msg)
	assert.Nil(t, err, "SetCurrency error")
}
