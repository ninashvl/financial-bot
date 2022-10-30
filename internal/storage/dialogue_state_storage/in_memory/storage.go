package in_memory

import (
	"sync"
)

type Storage struct {
	m     map[int64]int
	mutex sync.RWMutex
}

func New() *Storage {
	return &Storage{
		m:     make(map[int64]int),
		mutex: sync.RWMutex{},
	}
}

func (s *Storage) Add(userID int64, state int) {
	s.mutex.Lock()
	s.m[userID] = state
	s.mutex.Unlock()
}

func (s *Storage) DeleteState(userId int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.m, userId)
}

func (s *Storage) Get(userID int64) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.m[userID]
}
