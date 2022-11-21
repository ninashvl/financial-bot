package expense_storage

import (
	"context"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
)

//go:generate mockgen -source=istorage.go -destination=./mocks/mocks.go -package=mocks IStorage

const (
	Day = iota
	Month
	Year
)

type IStorage interface {
	Add(ctx context.Context, userID int64, expense *models.Expense) error
	GetByRange(ctx context.Context, userID int64, timeRange int) ([]*models.TotalExpense, error)
	SetCurrency(ctx context.Context, userID int64, curr string) error
	GetCurrency(ctx context.Context, userID int64) (string, error)
	UpdateCurrency(ctx context.Context) error
	SetLimit(ctx context.Context, userID int64, limit float64) error
	GetLimit(ctx context.Context, userID int64) (float64, error)
}

type IExpenseCache interface {
	GetByRange(ctx context.Context, userID int64, timeRange int) ([]*models.TotalExpense, error)
	SetByRange(ctx context.Context, userID int64, timeRange int, exps []*models.TotalExpense) error
}
