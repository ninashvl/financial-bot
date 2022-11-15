package psql

import (
	"context"
	"database/sql"
	"log"
	"time"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/clients/tradingview"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
)

var _ expense_storage.IStorage = &Storage{}

type Storage struct {
	p *sql.DB

	currencyClient tradingview.Client

	usdRUB float64
	cnyRUB float64
	eurRUB float64
}

func New(pool *sql.DB) *Storage {
	return &Storage{p: pool}
}

func (s *Storage) Add(ctx context.Context, userID int64, expense *models.Expense) error {
	_, err := s.p.ExecContext(ctx,
		"INSERT INTO expenses (user_id, amount, exp_category, exp_date) VALUES ($1, $2, $3, $4)",
		userID, expense.Amount, expense.Category, expense.Date)
	return err
}

func (s *Storage) GetByRange(ctx context.Context, userID int64, timeRange int) ([]*models.TotalExpense, error) {
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
	rows, err := s.p.QueryContext(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*models.TotalExpense
	for rows.Next() {
		t := &models.TotalExpense{}
		err = rows.Scan(&t.Category, &t.Amount)
		if err != nil {
			return nil, err
		}
		res = append(res, t)

	}

	curr, err := s.GetCurrency(ctx, userID)
	if err != nil {
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

	return res, nil
}

func (s *Storage) SetCurrency(ctx context.Context, userID int64, curr string) error {
	_, err := s.p.ExecContext(ctx, "INSERT INTO user_currency (user_id, currency_value) VALUES ($1, $2) ON CONFLICT (user_id) "+
		"DO UPDATE SET currency_value = $2 WHERE user_currency.user_id = $1", userID, curr)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetCurrency(ctx context.Context, userID int64) (string, error) {
	var res string
	err := s.p.QueryRowContext(ctx, "SELECT currency_value FROM user_currency WHERE user_id = $1", userID).Scan(&res)
	if err != nil {
		return "", err
	}
	if res == "" {
		return models.RubCurrency, nil
	}
	return res, nil
}

func (s *Storage) UpdateCurrency(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute * 10)
	if err := s.updatingCurrency(ctx); err != nil {
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
	log.Println("[INFO] start of updating currency quotes")

	curr, err := s.currencyClient.GetQuote(tradingview.UsdTicker)
	if err != nil {
		log.Println("[ERROR] getting usd quote error", err.Error())
		return err
	}
	s.usdRUB = curr
	curr, err = s.currencyClient.GetQuote(tradingview.EurTicker)
	if err != nil {
		log.Println("[ERROR] getting eur quote error", err.Error())
		return err
	}
	s.eurRUB = curr
	curr, err = s.currencyClient.GetQuote(tradingview.CnyTicker)
	if err != nil {
		log.Println("[ERROR] getting cny quote error", err.Error())
		return err
	}
	s.cnyRUB = curr
	log.Println("[INFO] usd quote =", s.usdRUB)
	log.Println("[INFO] eur quote =", s.eurRUB)
	log.Println("[INFO] cny quote =", s.cnyRUB)
	return nil
}

func (s *Storage) SetLimit(ctx context.Context, userID int64, limit float64) error {
	_, err := s.p.ExecContext(ctx, "INSERT INTO user_limit (user_id, limit_value) VALUES ($1, $2) ON CONFLICT (user_id) "+
		"DO UPDATE SET limit_value = $2 WHERE user_limit.user_id = $1", userID, limit)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetLimit(ctx context.Context, userID int64) (float64, error) {
	var res float64
	err := s.p.QueryRowContext(ctx, "SELECT limit_value FROM user_limit WHERE user_id = $1", userID).Scan(&res)
	if err != nil {
		return 0, err
	}

	curr, err := s.GetCurrency(ctx, userID)
	if err != nil {
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
