package expense_storage

import (
	"context"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
)

const (
	Day = iota
	Month
	Year
)

type IStorage interface {
	Add(userID int64, expense *models.Expense)
	GetByRange(userID int64, timeRange int) ([]*models.TotalExpense, error)
	SetCurrency(userID int64, curr string)
	GetCurrency(userID int64) string
	UpdateCurrency(ctx context.Context)
}
