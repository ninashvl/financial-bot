package models

import (
	"time"
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
