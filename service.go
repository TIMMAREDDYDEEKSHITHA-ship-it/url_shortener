package main

import (
	"context"
	"errors"
)

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	GetUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *User) error {
	if user.ID == "" || user.Email == "" {
		return errors.New("id and email are required")
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUsers(ctx context.Context) ([]User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) UpdateUser(ctx context.Context, user *User) error {
	if user.ID == "" || user.Email == "" {
		return errors.New("id and email are required")
	}
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(ctx, id)
}

type mockService struct{}

func (m *mockService) CreateUser(ctx context.Context, user *User) error {
	return nil
}
func (m *mockService) GetUsers(ctx context.Context) ([]User, error) {
	return []User{
		{ID: "1", Name: "TestUser", Email: "test@example.com"},
	}, nil
}
func (m *mockService) UpdateUser(ctx context.Context, user *User) error {
	if user.ID == "999" {
		return errors.New("not found")
	}
	return nil
}
func (m *mockService) DeleteUser(ctx context.Context, id string) error {
	if id == "999" {
		return errors.New("not found")
	}
	return nil
}
