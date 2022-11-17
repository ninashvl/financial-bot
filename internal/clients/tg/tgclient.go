package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/messages"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client *tgbotapi.BotAPI
	logger zerolog.Logger
}

func New(tokenGetter TokenGetter, l zerolog.Logger) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client: client,
		logger: l,
	}, nil
}

func (c *Client) SendMessage(text string, userID int64) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := c.client.Send(msg)
	if err != nil {
		c.logger.Error().Err(err)
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(ctx context.Context, bot *messages.Bot) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)
	go bot.ListenQuotes(ctx)
	c.logger.Info().Msg("listening for messages")
	for {
		select {
		case <-ctx.Done():
			c.logger.Info().Msg("Stop listening for messages")
			return
		case update := <-updates:
			if update.Message != nil { // If we got a message
				c.logger.Info().Str("message", update.Message.From.UserName).Str("text", update.Message.Text)
				msg := &messages.Message{
					Text:      update.Message.Text,
					UserID:    update.Message.From.ID,
					IsCommand: update.Message.IsCommand(),
				}
				err := bot.IncomingMessage(ctx, msg)
				if err != nil {
					c.logger.Error().Err(err).Msg("error processing message:")
				}
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

func (c *Client) SendCurrencyKeyboard(userID int64, text string) error {
	rangeKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(models.UsdCurrency),
			tgbotapi.NewKeyboardButton(models.RubCurrency),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(models.EurCurrency),
			tgbotapi.NewKeyboardButton(models.CnyCurrency),
		),
	)
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyMarkup = rangeKeyboard
	_, err := c.client.Send(msg)
	return err
}
