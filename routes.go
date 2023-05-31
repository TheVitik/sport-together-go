package main

import (
	"fmt"
	"github.com/TheVitik/sport-together-go/internal/handlers"
	"github.com/TheVitik/sport-together-go/internal/middlewares"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var Router *mux.Router

func InitRoutes(handler *handlers.Handler) {
	Router = mux.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	Router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintln(writer, "HELLO")
		if err != nil {
			return
		}
	})
	// Event endpoints
	Router.HandleFunc("/events", middlewares.Auth(handler.CreateEvent)).Methods("POST")
	Router.HandleFunc("/events/{id}", handler.GetEvent).Methods("GET")
	Router.HandleFunc("/events", handler.GetAllEvents).Methods("GET")
	Router.HandleFunc("/events/{id}", middlewares.Auth(handler.UpdateEvent)).Methods("PUT")
	Router.HandleFunc("/events/{id}", middlewares.Auth(handler.DeleteEvent)).Methods("DELETE")

	// User endpoints
	Router.HandleFunc("/register", handler.Register).Methods("POST")
	Router.HandleFunc("/login", handler.Login).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+port, Router))
}
