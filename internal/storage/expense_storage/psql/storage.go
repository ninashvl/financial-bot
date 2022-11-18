package psql

import (
	"context"
	sqllib "database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/clients/tradingview"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ expense_storage.IStorage = &Storage{}

type Storage struct {
	p              *sqllib.DB
	logger         zerolog.Logger
	currencyClient tradingview.Client

	usdRUB float64
	cnyRUB float64
	eurRUB float64
}

func New(pool *sqllib.DB, l zerolog.Logger) *Storage {
	return &Storage{p: pool, logger: l}
}

func (s *Storage) Add(ctx context.Context, userID int64, expense *models.Expense) error {

	var span trace.Span
	var err error
	ctx, span = otel.Tracer("update").Start(ctx, "storage.Add")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	_, err = s.p.ExecContext(ctx,
		"INSERT INTO expenses (user_id, amount, exp_category, exp_date) VALUES ($1, $2, $3, $4)",
		userID, expense.Amount, expense.Category, expense.Date)
	s.logger.Info().Int64("user", userID).Float64("expense", expense.Amount).Str("category", expense.Category).Err(err).Msg("expense added to storage")
	return err

}

func (s *Storage) GetByRange(ctx context.Context, userID int64, timeRange int) ([]*models.TotalExpense, error) {
	var span trace.Span
	var err error
	ctx, span = otel.Tracer("update").Start(ctx, "storage.GetByRange")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	sql := ""
	switch timeRange {
	case expense_storage.Day:
		sql = "SELECT exp_category, SUM(amount) FROM expenses " +
			"WHERE user_id = $1 AND date_part('day',exp_date)= date_part('day', CURRENT_DATE) AND " +
			"date_part('month',exp_date)= date_part('month', CURRENT_DATE) and " +
			"date_part('year',exp_date)= date_part('year', CURRENT_DATE) GROUP BY exp_category"
	case expense_storage.Month:
		sql = "SELECT exp_category, SUM(amount) FROM expenses " +
			"WHERE user_id = $1 AND date_part('month',exp_date)= date_part('month', CURRENT_DATE) AND " +
			"date_part('year',exp_date)= date_part('year', CURRENT_DATE) GROUP BY exp_category"

	case expense_storage.Year:
		sql = "SELECT exp_category, SUM(amount) FROM expenses " +
			"WHERE user_id = $1 AND " +
			"date_part('year',exp_date)= date_part('year', CURRENT_DATE) GROUP BY exp_category"
	}
	var rows *sqllib.Rows
	rows, err = s.p.QueryContext(ctx, sql, userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("query context error")
		return nil, err
	}
	defer rows.Close()
	var res []*models.TotalExpense
	for rows.Next() {
		t := &models.TotalExpense{}
		err = rows.Scan(&t.Category, &t.Amount)
		if err != nil {
			s.logger.Error().Err(err).Msg("rows scan error")
			return nil, err
		}
		res = append(res, t)
	}

	curr, err := s.GetCurrency(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("get currency error")
		return nil, err
	}
	var quoteVal float64 = 1
	switch curr {
	case models.UsdCurrency:
		quoteVal = s.usdRUB
	case models.CnyCurrency:
		quoteVal = s.cnyRUB
	case models.EurCurrency:
		quoteVal = s.eurRUB
	}

	for _, expense := range res {
		expense.Amount /= quoteVal
	}
	s.logger.Info().Int64("user", userID).Int("timerange", timeRange).Msg("get by range func executed")
	return res, nil
}

func (s *Storage) SetCurrency(ctx context.Context, userID int64, curr string) error {
	var span trace.Span
	var err error
	ctx, span = otel.Tracer("update").Start(ctx, "storage.SetCurrency")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	_, err = s.p.ExecContext(ctx, "INSERT INTO user_currency (user_id, currency_value) VALUES ($1, $2) ON CONFLICT (user_id) "+
		"DO UPDATE SET currency_value = $2 WHERE user_currency.user_id = $1", userID, curr)
	if err != nil {
		s.logger.Error().Err(err).Msg("exec context error")
		return err
	}
	return nil
}

func (s *Storage) GetCurrency(ctx context.Context, userID int64) (string, error) {
	var span trace.Span
	var err error
	ctx, span = otel.Tracer("update").Start(ctx, "storage.GetCurrency")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var res string
	err = s.p.QueryRowContext(ctx, "SELECT currency_value FROM user_currency WHERE user_id = $1", userID).Scan(&res)
	if err != nil {
		if errors.Is(err, sqllib.ErrNoRows) {
			return models.RubCurrency, nil
		}
		s.logger.Error().Err(err)
		return "", err
	}
	if res == "" {
		return models.RubCurrency, nil
	}
	return res, nil
}

func (s *Storage) UpdateCurrency(ctx context.Context) error {
	s.logger.Info().Msg("Currency update started")
	ticker := time.NewTicker(time.Minute * 10)
	if err := s.updatingCurrency(ctx); err != nil {
		s.logger.Error().Err(err).Msg("updating currency error")
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.updatingCurrency(ctx); err != nil {
				continue
			}
		}
	}
}

func (s *Storage) updatingCurrency(ctx context.Context) error {
	var span trace.Span
	var err error
	ctx, span = otel.Tracer("update_currency").Start(ctx, "updatingCurrency")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	s.logger.Info().Msg("Start of updating currency quotes")
	curr, err := s.currencyClient.GetQuote(ctx, tradingview.UsdTicker)
	if err != nil {
		s.logger.Error().Err(err).Msg("getting usd quote error")
		return err
	}
	s.usdRUB = curr
	curr, err = s.currencyClient.GetQuote(ctx, tradingview.EurTicker)
	if err != nil {
		s.logger.Error().Err(err).Msg("getting eur quote error")
		return err
	}
	s.eurRUB = curr
	curr, err = s.currencyClient.GetQuote(ctx, tradingview.CnyTicker)
	if err != nil {
		s.logger.Error().Err(err).Msg("getting cny quote error")
		return err
	}
	s.cnyRUB = curr
	s.logger.Info().Float64("usd", s.usdRUB).Float64("eur", s.eurRUB).Float64("cny", s.cnyRUB).Msg("current quotes")
	return nil
}

func (s *Storage) SetLimit(ctx context.Context, userID int64, limit float64) error {
	var span trace.Span
	var err error
	ctx, span = otel.Tracer("update").Start(ctx, "storage.SetLimit")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	_, err = s.p.ExecContext(ctx, "INSERT INTO user_limit (user_id, limit_value) VALUES ($1, $2) ON CONFLICT (user_id) "+
		"DO UPDATE SET limit_value = $2 WHERE user_limit.user_id = $1", userID, limit)
	if err != nil {
		s.logger.Error().Int64("user", userID).Msg("Set limit error")
		return err
	}
	return nil
}

func (s *Storage) GetLimit(ctx context.Context, userID int64) (float64, error) {
	var span trace.Span
	var err error
	ctx, span = otel.Tracer("update").Start(ctx, "storage.GetLimit")
	defer func() {
		if err != nil && !errors.Is(err, sqllib.ErrNoRows) {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var res float64
	err = s.p.QueryRowContext(ctx, "SELECT limit_value FROM user_limit WHERE user_id = $1", userID).Scan(&res)
	if err != nil {
		s.logger.Error().Err(err).Msg("queryrowcontext error")
		return 0, err
	}

	curr, err := s.GetCurrency(ctx, userID)
	if err != nil {
		s.logger.Error().Int64("user", userID).Err(err).Msg("GetCurrency error")
		return 0, err
	}
	switch curr {
	case models.UsdCurrency:
		res /= s.usdRUB
	case models.CnyCurrency:
		res /= s.cnyRUB
	case models.EurCurrency:
		res /= s.eurRUB
	}
	return res, nil
}
