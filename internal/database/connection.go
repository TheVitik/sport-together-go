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
	//connStr := "postgresql://postgres:11111111@127.0.0.1/postgres?sslmode=disable"
	connStr := "postgres://ldbdfmntglzcku:a4b8bfea50dd08ff6296b7255318e34b6924f04e6645c4f48cc2ae2a85d47c94@ec2-34-236-103-63.compute-1.amazonaws.com:5432/ddh2kc55mhedcp"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}
	return &Connection{db}
}

func (c *Connection) Migrate() {
	_, err := c.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(64),
		email VARCHAR(64),
    	password VARCHAR(64)
	)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.Exec(`CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		name VARCHAR(64),
		date VARCHAR(32),
    	details TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}
