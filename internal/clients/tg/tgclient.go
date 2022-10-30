package tg

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/messages"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client *tgbotapi.BotAPI
}

func New(tokenGetter TokenGetter) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SendMessage(text string, userID int64) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(msgModel *messages.Bot) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	log.Println("listening for messages")

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := &messages.Message{
				Text:      update.Message.Text,
				UserID:    update.Message.From.ID,
				IsCommand: update.Message.IsCommand(),
			}
			err := msgModel.IncomingMessage(msg)
			if err != nil {
				log.Println("error processing message:", err)
			}
		}
	}
}

func (c *Client) SendRangeKeyboard(userID int64, text string) error {
	rangeKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("День"),
			tgbotapi.NewKeyboardButton("Месяц"),
			tgbotapi.NewKeyboardButton("Год"),
		),
	)
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyMarkup = rangeKeyboard
	_, err := c.client.Send(msg)
	return err
}
