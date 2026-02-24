package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func initDB() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/url_shortener?sslmode=disable"
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	log.Println("Database connected successfully")
}

func createTables(ctx context.Context) error {
	_, err := db.NewCreateTable().
		Model((*User)(nil)).
		IfNotExists().
		Exec(ctx)
	return err
}

func main() {
	initDB()                //intializing the DB connection
	server := NewServer(db) //create a new server instance with the DB connection	db dependency

	http.HandleFunc("/users", server.usersHandler) //pass the handler to the server method
	log.Println("Server running on port:8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) //starts  the http server

}
