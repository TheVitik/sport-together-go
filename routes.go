package main

import (
	"fmt"
	"github.com/TheVitik/sport-together-go/internal/handlers"
	"github.com/TheVitik/sport-together-go/internal/middlewares"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func initRoutes(handler *handlers.Handler) {
	router := mux.NewRouter()
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintln(writer, "HELLO")
		if err != nil {
			return
		}
	})
	// Event endpoints
	router.HandleFunc("/events", middlewares.Auth(handler.CreateEvent)).Methods("POST")
	router.HandleFunc("/events/{id}", handler.GetEvent).Methods("GET")
	router.HandleFunc("/events", handler.GetAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", middlewares.Auth(handler.UpdateEvent)).Methods("PUT")
	router.HandleFunc("/events/{id}", middlewares.Auth(handler.DeleteEvent)).Methods("DELETE")

	// User endpoints
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/login", handler.Login).Methods("POST")

	log.Fatal(http.ListenAndServe(":80", router))
}
