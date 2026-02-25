package main

import (
	"context"

	"github.com/uptrace/bun"
)

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepo {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.NewSelect().Model(&users).Scan(ctx)
	return users, err
}

// Add these two methods:
func (r *UserRepository) Update(ctx context.Context, user *User) error {
	_, err := r.db.NewUpdate().Model(user).Where("id = ?", user.ID).Exec(ctx)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model((*User)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
