package main

import (
	"github.com/TheVitik/sport-together-go/internal/database"
	"github.com/TheVitik/sport-together-go/internal/handlers"
	"github.com/TheVitik/sport-together-go/internal/repositories"
)

func main() {
	connection := database.NewConnection()
	connection.Migrate()

	repository := repositories.NewRepository(connection)
	handler := handlers.NewHandler(repository)

	initRoutes(handler)
}
