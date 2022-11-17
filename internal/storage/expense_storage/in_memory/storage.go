package in_memory

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/clients/tradingview"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
)

var _ expense_storage.IStorage = &Storage{}

type Storage struct {
	m              map[int64]map[string][]*models.Expense
	currency       map[int64]string
	limit          map[int64]float64
	currencyClient tradingview.Client

	usdRUB float64
	cnyRUB float64
	eurRUB float64

	mutex sync.RWMutex
}

func (s *Storage) GetCurrency(ctx context.Context, userID int64) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if _, ok := s.currency[userID]; !ok {
		return models.RubCurrency, nil
	}
	return s.currency[userID], nil
}

func (s *Storage) SetCurrency(ctx context.Context, userID int64, curr string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.currency[userID] = curr
	return nil
}

func New() *Storage {
	return &Storage{
		m:        make(map[int64]map[string][]*models.Expense),
		currency: make(map[int64]string),
		limit:    make(map[int64]float64),
		mutex:    sync.RWMutex{},
		usdRUB:   61,
		cnyRUB:   8,
		eurRUB:   61,
	}
}

func (s *Storage) Add(ctx context.Context, userID int64, expense *models.Expense) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.m[userID]; !ok {
		s.m[userID] = make(map[string][]*models.Expense)
	}
	if _, ok := s.m[userID][expense.Category]; !ok {
		s.m[userID][expense.Category] = make([]*models.Expense, 0)
	}
	if c, ok := s.currency[userID]; ok {
		switch c {
		case models.UsdCurrency:
			expense.Amount = expense.Amount * s.usdRUB
		case models.CnyCurrency:
			expense.Amount = expense.Amount * s.cnyRUB
		case models.EurCurrency:
			expense.Amount = expense.Amount * s.eurRUB
		}
	}
	s.m[userID][expense.Category] = append(s.m[userID][expense.Category], expense)
	return nil
}

func (s *Storage) GetByRange(ctx context.Context, userID int64, timeRange int) ([]*models.TotalExpense, error) {
	m := make(map[string]*models.TotalExpense)
	now := time.Now()
	if _, ok := s.m[userID]; !ok {
		return nil, errors.New("Пока у Вас еще не было трат за выбранный диапазон...")
	}
	var curr float64 = 1
	if c, ok := s.currency[userID]; ok {
		switch c {
		case models.UsdCurrency:
			curr = s.usdRUB
		case models.CnyCurrency:
			curr = s.cnyRUB
		case models.EurCurrency:
			curr = s.eurRUB
		}
	}
	for category, exps := range s.m[userID] {
		m[category] = &models.TotalExpense{}
		m[category].Category = category
		for _, v := range exps {
			switch timeRange {
			case expense_storage.Day:
				if v.Date.Day() == now.Day() && v.Date.Month() == now.Month() && v.Date.Year() == now.Year() {
					m[category].Amount += v.Amount / curr
				}
			case expense_storage.Month:
				if v.Date.Month() == now.Month() && v.Date.Year() == now.Year() {
					m[category].Amount += v.Amount / curr
				}
			case expense_storage.Year:
				if v.Date.Year() == now.Year() {
					m[category].Amount += v.Amount / curr
				}
			}

		}
	}
	res := make([]*models.TotalExpense, 0, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res, nil
}

// todo обновить логирование
func (s *Storage) UpdateCurrency(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			log.Println("[INFO] start of updating currency quotes")

			curr, err := s.currencyClient.GetQuote(tradingview.UsdTicker)
			if err != nil {
				log.Println("[ERROR] getting usd quote error", err.Error())

				continue
			}
			s.usdRUB = curr
			curr, err = s.currencyClient.GetQuote(tradingview.EurTicker)
			if err != nil {
				log.Println("[ERROR] getting eur quote error", err.Error())
				continue
			}
			s.eurRUB = curr
			curr, err = s.currencyClient.GetQuote(tradingview.CnyTicker)
			if err != nil {
				log.Println("[ERROR] getting cny quote error", err.Error())
				continue
			}
			s.cnyRUB = curr
			log.Println("[INFO] usd quote =", s.usdRUB)
			log.Println("[INFO] eur quote =", s.eurRUB)
			log.Println("[INFO] cny quote =", s.cnyRUB)

		}

	}
}

func (s *Storage) SetLimit(ctx context.Context, userID int64, limit float64) error {
	s.mutex.Lock()
	s.limit[userID] = limit
	s.mutex.Unlock()
	return nil
}

func (s *Storage) GetLimit(ctx context.Context, userID int64) (float64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.limit[userID], nil
}
