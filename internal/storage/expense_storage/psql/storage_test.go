package psql

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
)

func TestStorage_Add_and_GetByRange(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Creation sql mock error")

	ctx := context.Background()

	userID := int64(1)
	amount := 100.
	cat := "test"

	s := &Storage{
		p:      db,
		eurRUB: 63,
		usdRUB: 60,
		cnyRUB: 9,
	}

	// сегодня RUB
	sql1 := regexp.QuoteMeta("SELECT exp_category, SUM(amount) FROM expenses " +
		"WHERE user_id = $1 AND date_part('day',exp_date)= date_part('day', CURRENT_DATE) AND " +
		"date_part('month',exp_date)= date_part('month', CURRENT_DATE) and " +
		"date_part('year',exp_date)= date_part('year', CURRENT_DATE) GROUP BY exp_category")

	sql2 := regexp.QuoteMeta("SELECT currency_value FROM user_currency WHERE")

	mock.ExpectQuery(sql1).WithArgs(userID).WillReturnRows(mock.NewRows([]string{"exp_category", "sum"}).AddRow(cat, amount))

	mock.ExpectQuery(sql2).WithArgs(userID).WithArgs().WillReturnRows(mock.NewRows([]string{"currency_value"}).AddRow("RUB"))

	res, err := s.GetByRange(ctx, userID, expense_storage.Day)
	assert.Nil(t, err, "GetByRange error")

	assert.Equal(t, 1, len(res))
	assert.Equal(t, amount, res[0].Amount)
	assert.Equal(t, cat, res[0].Category)

	// сегодня EUR
	mock.ExpectQuery(sql1).WithArgs(userID).WillReturnRows(mock.NewRows([]string{"exp_category", "sum"}).AddRow(cat, amount))

	mock.ExpectQuery(sql2).WithArgs(userID).WithArgs().WillReturnRows(mock.NewRows([]string{"currency_value"}).AddRow("EUR"))

	res, err = s.GetByRange(ctx, userID, expense_storage.Day)
	assert.Nil(t, err, "GetByRange error")

	assert.Equal(t, 1, len(res))
	assert.Equal(t, amount, res[0].Amount*s.eurRUB)
	assert.Equal(t, cat, res[0].Category)

	assert.Nil(t, mock.ExpectationsWereMet(), "Not all mocks are executed")
}

func TestStorage_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Creation sql mock error")
	ctx := context.Background()

	userID := int64(1)
	amount := 100.
	cat := "test"

	s := &Storage{
		p:      db,
		eurRUB: 63,
		usdRUB: 60,
		cnyRUB: 9,
	}
	sql := regexp.QuoteMeta("INSERT INTO expenses (user_id, amount, exp_category, exp_date)")

	now := time.Now()
	mock.ExpectExec(sql).WithArgs(userID, amount, cat, now).WillReturnResult(sqlmock.NewResult(1, 1))

	yesterday := time.Now().Add(-time.Hour * 24)
	mock.ExpectExec(sql).WithArgs(userID, amount, cat, yesterday).WillReturnResult(sqlmock.NewResult(1, 1))

	lastMonth := time.Now().Add(-time.Hour * 24 * 32)
	mock.ExpectExec(sql).WithArgs(userID, amount, cat, lastMonth).WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.Add(ctx, userID, &models.Expense{
		Amount:   amount,
		Category: cat,
		Date:     now,
	})

	assert.Nil(t, err, "Add expense error")

	err = s.Add(ctx, userID, &models.Expense{
		Amount:   amount,
		Category: cat,
		Date:     yesterday,
	})

	assert.Nil(t, err, "Add expense error")

	err = s.Add(ctx, userID, &models.Expense{
		Amount:   amount,
		Category: cat,
		Date:     lastMonth,
	})

	assert.Nil(t, err, "Add expense error")

	assert.Nil(t, mock.ExpectationsWereMet(), "Not all mocks are executed")
}
