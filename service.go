package main

type Service struct {
	counterManager *CounterManager
}

func NewService() *Service {
	return &Service{
		counterManager: NewCounterManager(),
	}
}

func (s *Service) Increment(key string) {
	s.counterManager.Increment(key)
}

func (s *Service) GetCounter(key string) int64 {
	return s.counterManager.Get(key)
}
