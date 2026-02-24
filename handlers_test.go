package main

import (
	"bytes"
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestDB(t *testing.T) {
	dsn := "postgres://timmareddydeekshitha@localhost:5432/url_shortener?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())
	err := db.Ping()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// adding below code to ensure the table exists
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
	setupTestDB(t)
	server := NewServer(db) //create server instance
	reqBody := []byte(`{"id":"1","name":"Test","email":"test@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.usersHandler(w, req) //call the handler
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}
}

func TestGetUsers(t *testing.T) {
	setupTestDB(t)
	server := NewServer(db)
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	server.usersHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	setupTestDB(t)
	server := NewServer(db)
	user := User{ID: "2", Name: "Old", Email: "old@example.com"}
	db.NewInsert().Model(&user).Exec(context.Background())
	reqBody := []byte(`{"name":"New","email":"new@example.com"}`)
	req := httptest.NewRequest(http.MethodPut, "/users?id=2", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.usersHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	setupTestDB(t)
	server := NewServer(db)
	user := User{ID: "3", Name: "DeleteMe", Email: "delete@example.com"}
	db.NewInsert().Model(&user).Exec(context.Background())
	req := httptest.NewRequest(http.MethodDelete, "/users?id=3", nil)
	w := httptest.NewRecorder()
	server.usersHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestInvalidJSON(t *testing.T) {
	setupTestDB(t)
	server := NewServer(db)
	reqBody := []byte(`invalid json`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	server.usersHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	setupTestDB(t)
	server := NewServer(db)
	req := httptest.NewRequest(http.MethodDelete, "/users?id=999", nil)
	w := httptest.NewRecorder()
	server.usersHandler(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	setupTestDB(t)
	server := NewServer(db)
	reqBody := []byte(`{"name":"New","email":"new@example.com"}`)
	req := httptest.NewRequest(http.MethodPut, "/users?id=999", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.usersHandler(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestMissingIDParameter(t *testing.T) {
	setupTestDB(t)
	server := NewServer(db)
	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	w := httptest.NewRecorder()
	server.usersHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	setupTestDB(t)
	server := NewServer(db)
	req := httptest.NewRequest(http.MethodPatch, "/users", nil)
	w := httptest.NewRecorder()
	server.usersHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)

	}
}
