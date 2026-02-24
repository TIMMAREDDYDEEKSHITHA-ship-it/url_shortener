package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type mockRepo struct{}

func (m *mockRepo) Create(ctx context.Context, user *User) error {
	return nil
}

func (m *mockRepo) GetAll(ctx context.Context) ([]User, error) {
	return []User{
		{ID: "1", Name: "TestUser", Email: "test@example.com"},
	}, nil
}

func (m *mockRepo) Update(ctx context.Context, user *User) error {
	if user.ID == "999" {
		return errors.New("not found")
	}
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	if id == "999" {
		return errors.New("not found")
	}
	return nil
}

func setupTestDB(t *testing.T) {
	dsn := "postgres://timmareddydeekshitha@localhost:5432/url_shortener?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())
	err := db.Ping()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()
	_, err = db.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	_, err = db.NewTruncateTable().Model((*User)(nil)).Exec(ctx)
	if err != nil {
		t.Fatalf("Failed to truncate users table:%v", err)
	}
}

func TestCreateUser(t *testing.T) {
	repo := &mockRepo{}
	handler := NewHandler(repo)
	reqBody := []byte(`{"id":"100","name":"Test","email":"test@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Users(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestGetUsers(t *testing.T) {
	repo := &mockRepo{}
	handler := NewHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	handler.Users(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	setupTestDB(t)
	repo := NewUserRepository(db)
	handler := NewHandler(repo)
	user := User{ID: "2", Name: "Old", Email: "old@example.com"}
	db.NewInsert().Model(&user).Exec(context.Background())
	reqBody := []byte(`{"name":"New","email":"new@example.com"}`)
	req := httptest.NewRequest(http.MethodPut, "/users?id=2", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Users(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	setupTestDB(t)
	repo := NewUserRepository(db)
	handler := NewHandler(repo)
	user := User{ID: "3", Name: "DeleteMe", Email: "delete@example.com"}
	db.NewInsert().Model(&user).Exec(context.Background())
	req := httptest.NewRequest(http.MethodDelete, "/users?id=3", nil)
	w := httptest.NewRecorder()
	handler.Users(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestInvalidJSON(t *testing.T) {
	setupTestDB(t)
	repo := NewUserRepository(db)
	handler := NewHandler(repo)
	reqBody := []byte(`invalid json`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	handler.Users(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	repo := &mockRepo{}
	handler := NewHandler(repo)
	req := httptest.NewRequest(http.MethodDelete, "/users?id=999", nil)
	w := httptest.NewRecorder()
	handler.Users(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	repo := &mockRepo{}
	handler := NewHandler(repo)
	reqBody := []byte(`{"name":"New","email":"new@example.com"}`)
	req := httptest.NewRequest(http.MethodPut, "/users?id=999", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Users(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestMissingIDParameter(t *testing.T) {
	setupTestDB(t)
	repo := NewUserRepository(db)
	handler := NewHandler(repo)
	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	w := httptest.NewRecorder()
	handler.Users(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	setupTestDB(t)
	repo := NewUserRepository(db)
	handler := NewHandler(repo)
	req := httptest.NewRequest(http.MethodPatch, "/users", nil)
	w := httptest.NewRecorder()
	handler.Users(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
