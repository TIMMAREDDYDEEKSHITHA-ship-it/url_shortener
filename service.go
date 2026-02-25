package main

import (
	"context"
	"errors"
)

type URLService interface {
	CreateURL(ctx context.Context, url *URL) error
	GetByCode(ctx context.Context, code string) (*URL, error)
}

type urlService struct {
	dbRepo URLRepository
}

func NewURLService(repo URLRepository) URLService {
	return &urlService{
		dbRepo: repo,
	}
}

func (s *urlService) CreateURL(ctx context.Context, url *URL) error {
	if url.Code == "" || url.LongURL == "" {
		return errors.New("code and long_url are required")
	}
	return s.dbRepo.Create(ctx, url)
}

func (s *urlService) GetByCode(ctx context.Context, code string) (*URL, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}
	return s.dbRepo.GetByCode(ctx, code)
}
