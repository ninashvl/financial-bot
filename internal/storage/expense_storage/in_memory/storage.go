package in_memory

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"gitlab.ozon.dev/ninashvl/homework-1/internal/models"
	"gitlab.ozon.dev/ninashvl/homework-1/internal/storage/expense_storage"
)

var _ expense_storage.IStorage = &Storage{}

type Storage struct {
	m     map[int64]map[string][]*models.Expense
	mutex sync.RWMutex
}

func New() *Storage {
	return &Storage{
		m:     make(map[int64]map[string][]*models.Expense),
		mutex: sync.RWMutex{},
	}
}

func (s *Storage) Add(userID int64, expense *models.Expense) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.m[userID]; !ok {
		s.m[userID] = make(map[string][]*models.Expense)
	}
	if _, ok := s.m[userID][expense.Category]; !ok {
		s.m[userID][expense.Category] = make([]*models.Expense, 0)
	}
	s.m[userID][expense.Category] = append(s.m[userID][expense.Category], expense)
}

func (s *Storage) GetByRange(userID int64, timeRange int) ([]*models.TotalExpense, error) {
	m := make(map[string]*models.TotalExpense)
	now := time.Now()
	if _, ok := s.m[userID]; !ok {
		return nil, errors.New("User not found")
	}
	for category, exps := range s.m[userID] {
		m[category] = &models.TotalExpense{}
		m[category].Category = category
		for _, v := range exps {
			switch timeRange {
			case expense_storage.Day:
				if v.Date.Day() == now.Day() && v.Date.Month() == now.Month() && v.Date.Year() == now.Year() {
					m[category].Amount += v.Amount
				}
			case expense_storage.Month:
				if v.Date.Month() == now.Month() && v.Date.Year() == now.Year() {
					m[category].Amount += v.Amount
				}
			case expense_storage.Year:
				if v.Date.Year() == now.Year() {
					m[category].Amount += v.Amount
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
