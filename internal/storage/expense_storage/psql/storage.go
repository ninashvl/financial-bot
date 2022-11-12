package psql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
)

var _ expense_storage.IStorage = &Storage{}

type Storage struct {
	p *pgxpool.Pool
}

func (s Storage) Add(ctx context.Context, userID int64, expense *models.Expense) error {
	_, err := s.p.Exec(context.TODO(),
		"INSERT INTO expenses (user_id, amount, exp_category, exp_date) VALUES ($1, $2, $3, $4)",
		userID, expense.Amount, expense.Category, expense.Date)
	return err
}

func (s Storage) GetByRange(ctx context.Context, userID int64, timeRange int) ([]*models.TotalExpense, error) {
	sql := ""
	switch timeRange {
	case expense_storage.Day:
		sql = "SELECT exp_category, amount FROM expeses " +
			"WHERE user_id = $1 AND date_part('day',exp_date)= date_part('day', CURRENT_DATE) AND " +
			"date_part('month',exp_date)= date_part('month', CURRENT_DATE) and " +
			"date_part('year',exp_date)= date_part('year', CURRENT_DATE) GROUP BY p_category"
	case expense_storage.Month:
		sql = "SELECT exp_category, amount FROM expeses " +
			"WHERE user_id = $1 AND date_part('month',exp_date)= date_part('month', CURRENT_DATE) AND " +
			"date_part('year',exp_date)= date_part('year', CURRENT_DATE) GROUP BY p_category"

	case expense_storage.Year:
		sql = "SELECT exp_category, amount FROM expeses " +
			"WHERE user_id = $1 AND " +
			"date_part('year',exp_date)= date_part('year', CURRENT_DATE) GROUP BY p_category"
	}
	rows, err := s.p.Query(ctx, sql, userID)
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
	return res, nil
}

func (s Storage) SetCurrency(ctx context.Context, userID int64, curr string) error {
	_, err := s.p.Exec(ctx, "INSERT INTO user_currency (user_id, currency_value) VALUES ($1, $2) ON CONFLICT (user_id) "+
		"DO UPDATE SET currency_value = $2 WHERE user_id = $1", userID, curr)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) GetCurrency(ctx context.Context, userID int64) (string, error) {
	var res string
	err := s.p.QueryRow(ctx, "SELECT currency_value FROM user_currency WHERE user_id = $1", userID).Scan(&res)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s Storage) UpdateCurrency(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{p: pool}
}
