package repositories

import "github.com/TheVitik/sport-together-go/internal/database"

type Repository struct {
	db *database.Connection
}

func NewRepository(db *database.Connection) *Repository {
	return &Repository{db: db}
}
