package main

import "github.com/uptrace/bun"

type Server struct {
	userRepo UserRepo
}

func NewServer(db *bun.DB) *Server {
	userRepo := NewUserRepository(db)
	return &Server{
		userRepo: userRepo,
	}
}
