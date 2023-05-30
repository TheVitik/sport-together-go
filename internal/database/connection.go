package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Connection struct {
	*sql.DB
}

func NewConnection() *Connection {
	connStr := "postgresql://postgres:11111111@127.0.0.1/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}
	return &Connection{db}
}
