package main

import (
	"context"

	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.NewSelect().Model(&users).Scan(ctx)
	return users, err
}
