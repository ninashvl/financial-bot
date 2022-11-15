package messages

import (
	"context"
	"strconv"
)

func (s *Bot) AddLimit(ctx context.Context, msg *Message) error {
	limit, err := strconv.ParseFloat(msg.Text, 64)
	if err != nil {
		return s.tgClient.SendMessage(invalidMsg, msg.UserID)
	}

	err = s.expStorage.SetLimit(ctx, msg.UserID, limit)
	if err != nil {
		return s.tgClient.SendMessage(err.Error(), msg.UserID)
	}
	return s.tgClient.SendMessage(limitSuccessfulSet, msg.UserID)
}
