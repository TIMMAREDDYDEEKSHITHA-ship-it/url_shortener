package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func setupTestServer(t *testing.T) *Server {
	dsn := "postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable"

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	testDB := bun.NewDB(sqldb, pgdialect.New())

	if err := testDB.Ping(); err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	// create table
	_, err := testDB.NewCreateTable().
		Model((*URL)(nil)).
		IfNotExists().
		Exec(context.Background())

	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	server := NewServer(testDB)
	server.RegisterRoutes()

	return server
}

func TestHealth(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestShorten(t *testing.T) {
	server := setupTestServer(t)

	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Invalid JSON response")
	}

	if resp["short_url"] == "" {
		t.Errorf("Expected short_url in response")
	}
}

func TestShortenInvalidJSON(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer([]byte(`invalid`)))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestRedirect(t *testing.T) {
	server := setupTestServer(t)

	// insert test URL manually
	url := &URL{
		Code:    "abc123",
		LongURL: "https://example.com",
	}

	_, err := server.db.NewInsert().Model(url).Exec(context.Background())
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected 302 redirect, got %d", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", w.Code)
	}
}
