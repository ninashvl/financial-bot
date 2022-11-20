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

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	expStore := expstorage.NewMockIStorage(ctrl)
	stStore := statestorage.NewMockIStorage(ctrl)

	bot := &Bot{
		tgClient:        sender,
		expStorage:      expStore,
		dlgStateStorage: stStore,
	}

	sender.EXPECT().SendMessage(gomock.Any(), "hello", int64(123))

	err := bot.IncomingMessage(context.TODO(), &Message{
		Text:      "/start",
		UserID:    123,
		IsCommand: true,
	})

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	expStore := expstorage.NewMockIStorage(ctrl)
	stStore := statestorage.NewMockIStorage(ctrl)

	bot := &Bot{
		tgClient:        sender,
		expStorage:      expStore,
		dlgStateStorage: stStore,
	}

	sender.EXPECT().SendMessage(gomock.Any(), invalidCommand, int64(123))

	err := bot.IncomingMessage(context.TODO(), &Message{
		Text:      "some text",
		UserID:    123,
		IsCommand: true,
	})

	assert.NoError(t, err)
}
