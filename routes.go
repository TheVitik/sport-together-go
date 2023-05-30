package main

import (
	"github.com/TheVitik/sport-together-go/internal/database"
	"github.com/TheVitik/sport-together-go/internal/handlers"
	"github.com/TheVitik/sport-together-go/internal/middlewares"
	"github.com/TheVitik/sport-together-go/internal/repositories"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func initRoutes() {
	router := mux.NewRouter()

	repository := repositories.NewRepository(database.NewConnection())
	handler := handlers.NewHandler(repository)

	// Event endpoints
	router.HandleFunc("/events", middlewares.Auth(handler.CreateEvent)).Methods("POST")
	router.HandleFunc("/events/{id}", handler.GetEvent).Methods("GET")
	router.HandleFunc("/events", handler.GetAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", middlewares.Auth(handler.UpdateEvent)).Methods("PUT")
	router.HandleFunc("/events/{id}", middlewares.Auth(handler.DeleteEvent)).Methods("DELETE")

	// User endpoints
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/login", handler.Login).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
