package main

import (
	"net/http"

	"github.com/uptrace/bun"
)

type Server struct {
	db      *bun.DB
	handler *Handler
	mux     *http.ServeMux
}

func NewServer(db *bun.DB) *Server {
	repo := NewURLRepository(db)
	service := NewURLService(repo)
	handler := NewHandler(service)

	return &Server{
		db:      db,
		handler: handler,
		mux:     http.NewServeMux(),
	}
}

func (s *Server) RegisterRoutes() {
	s.mux.HandleFunc("/health", s.handler.Health)
	s.mux.HandleFunc("/shorten", s.handler.Shorten)
	s.mux.HandleFunc("/", s.handler.Redirect)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
