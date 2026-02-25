package main

import (
	"sync"
)

type URLStore struct {
	mu   sync.RWMutex
	urls map[string]string // short code -> long URL
}

func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]string),
	}
}

func (s *URLStore) Save(shortCode, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[shortCode] = longURL
}

func (s *URLStore) Get(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, exists := s.urls[shortCode]
	return url, exists
}

func (s *URLStore) Delete(shortCode string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.urls, shortCode)
}
