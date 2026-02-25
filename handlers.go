package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type Handler struct {
	service URLService
}

func NewHandler(service URLService) *Handler {
	return &Handler{service: service}
}

func generateCode() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	code := generateCode()

	url := &URL{
		Code:    code,
		LongURL: req.URL,
	}

	if err := h.service.CreateURL(r.Context(), url); err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"short_url": "http://localhost:8080/" + code,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]

	if code == "" || code == "shorten" {
		http.NotFound(w, r)
		return
	}

	url, err := h.service.GetByCode(r.Context(), code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusFound)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
