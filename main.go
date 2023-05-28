package main

import (
	"database/sql"
	"github.com/TheVitik/sport-together-go/internal/app"
	ihttp "github.com/TheVitik/sport-together-go/internal/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	connStr := "postgresql://postgres:11111111@127.0.0.1/postgres?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	// Event endpoints
	router.HandleFunc("/events", ihttp.AuthMiddleware(app.CreateEvent)).Methods("POST")
	router.HandleFunc("/events/{id}", app.GetEvent).Methods("GET")
	router.HandleFunc("/events", app.GetAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", ihttp.AuthMiddleware(app.UpdateEvent)).Methods("PUT")
	router.HandleFunc("/events/{id}", ihttp.AuthMiddleware(app.DeleteEvent)).Methods("DELETE")

	// User endpoints
	router.HandleFunc("/register", app.Register(db)).Methods("POST")
	router.HandleFunc("/login", app.Login(db)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
