package models

import (
	"time"
)

const (
	RubCurrency = "RUB"
	UsdCurrency = "USD"
	EurCurrency = "EUR"
	CnyCurrency = "CNY"
)

type Expense struct {
	Amount   float64
	Category string
	Date     time.Time
}

type TotalExpense struct {
	Amount   float64
	Category string
}
