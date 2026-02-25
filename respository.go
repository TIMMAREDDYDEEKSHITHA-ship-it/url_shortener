package main

import (
	"context"

	"github.com/uptrace/bun"
)

type URLRepository interface {
	Create(ctx context.Context, url *URL) error
	GetByCode(ctx context.Context, code string) (*URL, error)
}

type urlRepository struct {
	db *bun.DB
}

func NewURLRepository(db *bun.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(ctx context.Context, url *URL) error {
	_, err := r.db.NewInsert().Model(url).Exec(ctx)
	return err
}

func (r *urlRepository) GetByCode(ctx context.Context, code string) (*URL, error) {
	var url URL

	err := r.db.NewSelect().
		Model(&url).
		Where("code = ?", code).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &url, nil
}
