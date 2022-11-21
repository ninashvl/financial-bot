package messages

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"gitlab.ozon.dev/ninashvl/homework-1/config"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/metrics"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/dialogue_state_storage"
	in_mem_dlg "gitlab.ozon.dev/ninashvl/homework-1/internal/storage/dialogue_state_storage/in_memory"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage/psql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type MessageSender interface {
	SendMessage(ctx context.Context, text string, userID int64) error
	SendRangeKeyboard(ctx context.Context, userID int64, text string) error
	SendCurrencyKeyboard(ctx context.Context, userID int64, text string) error
}

type Bot struct {
	tgClient        MessageSender
	expStorage      expense_storage.IStorage
	dlgStateStorage dialogue_state_storage.IStorage
	logger          zerolog.Logger
}

func New(tgClient MessageSender, cfg *config.Service, l zerolog.Logger) *Bot {
	db, err := sql.Open("pgx", cfg.PsqlDSN())
	if err != nil {
		l.Fatal().Err(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Bot{
		tgClient:        tgClient,
		expStorage:      psql.New(db, l),
		dlgStateStorage: in_mem_dlg.New(),
		logger:          l,
	}
}

type Message struct {
	Text      string
	UserID    int64
	IsCommand bool
}

func (s *Bot) HandleCommand(ctx context.Context, msg *Message) error {
	var span trace.Span
	ctx, span = otel.Tracer("update").Start(ctx, "bot.HandleCommand")
	defer span.End()
	metrics.TotalCommandsCount.WithLabelValues(msg.Text).Inc()
	s.logger.Debug().Str("text", msg.Text).Int64("user", msg.UserID).Msg("Handle command func called")

	switch {
	case msg.Text == startCommand:
		return s.tgClient.SendMessage(ctx, "hello", msg.UserID)
	case msg.Text == helpCommand && msg.IsCommand:
		return s.tgClient.SendMessage(ctx, help, msg.UserID)
	case msg.Text == addCommand && msg.IsCommand:
		s.dlgStateStorage.Set(msg.UserID, models.AddCommandState)
		return s.tgClient.SendMessage(ctx, addMessage, msg.UserID)
	case msg.Text == getExpensesCommand && msg.IsCommand:
		s.dlgStateStorage.Set(msg.UserID, models.GetCommandState)
		return s.tgClient.SendRangeKeyboard(ctx, msg.UserID, "Выберите диапазон")
	case msg.Text == chooseCurrencyCommand && msg.IsCommand:
		s.dlgStateStorage.Set(msg.UserID, models.ChooseCurrencyState)
		return s.tgClient.SendCurrencyKeyboard(ctx, msg.UserID, "Выберите валюту")
	case msg.Text == addLimit && msg.IsCommand:
		s.dlgStateStorage.Set(msg.UserID, models.AddLimit)
		return s.tgClient.SendMessage(ctx, "Введите лимит на траты в рублях на месяц", msg.UserID)
	default:
		return s.tgClient.SendMessage(ctx, invalidCommand, msg.UserID)
	}
}

func (s *Bot) HandleMessage(ctx context.Context, msg *Message) error {
	var span trace.Span
	ctx, span = otel.Tracer("update").Start(ctx, "bot.HandleMessage")
	defer span.End()

	s.logger.Debug().Str("text", msg.Text).Int64("user", msg.UserID).Msg("Handle message func called")
	switch {
	case !msg.IsCommand && s.dlgStateStorage.Get(msg.UserID) == models.AddCommandState:
		return s.AddExpense(ctx, msg)
	case !msg.IsCommand && s.dlgStateStorage.Get(msg.UserID) == models.GetCommandState:
		return s.GetExpense(ctx, msg)
	case !msg.IsCommand && s.dlgStateStorage.Get(msg.UserID) == models.ChooseCurrencyState:
		return s.SetCurrency(ctx, msg)
	case !msg.IsCommand && s.dlgStateStorage.Get(msg.UserID) == models.AddLimit:
		return s.AddLimit(ctx, msg)
	default:
		return s.tgClient.SendMessage(ctx, invalidMsg, msg.UserID)
	}
}

func (s *Bot) IncomingMessage(ctx context.Context, msg *Message) error {
	metrics.TotalMsgCount.Inc()
	defer func() {
		if !msg.IsCommand {
			s.dlgStateStorage.DeleteState(msg.UserID)
		}
	}()
	if msg.IsCommand {
		return s.HandleCommand(ctx, msg)
	}
	return s.HandleMessage(ctx, msg)
}

func (s *Bot) ListenQuotes(ctx context.Context) {
	err := s.expStorage.UpdateCurrency(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("UpdateCurrency failed")
	}
}
