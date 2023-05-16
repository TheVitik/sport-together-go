package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// User is a struct representing a user
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// users is a map of user IDs to User structs
var users = make(map[string]User)

// secretKey is the secret key used to sign JWT tokens
var secretKey = []byte("secret")

// Event is a struct representing an event
type Event struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Date    string `json:"date"`
	Details string `json:"details"`
}

// events is a slice of Event structs
var events []Event

// AuthMiddleware prevents unathorized access
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the request headers
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" format
		tokenString := strings.Split(authHeader, " ")[1]

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method and secret key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// CreateEvent creates a new event and adds it to the events slice
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	event.ID = len(events) + 1
	events = append(events, event)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// GetEvent returns an event with the specified ID
func GetEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, event := range events {
		if strconv.Itoa(event.ID) == params["id"] {
			json.NewEncoder(w).Encode(event)
			return
		}
	}
	http.Error(w, "Event not found", http.StatusNotFound)
}

// GetAllEvents returns all events in the events slice
func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

// UpdateEvent updates an event with the specified ID
func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, event := range events {
		if strconv.Itoa(event.ID) == params["id"] {
			events = append(events[:i], events[i+1:]...)
			var updatedEvent Event
			err := json.NewDecoder(r.Body).Decode(&updatedEvent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updatedEvent.ID = event.ID
			events = append(events, updatedEvent)
			json.NewEncoder(w).Encode(updatedEvent)
			return
		}
	}
	http.Error(w, "Event not found", http.StatusNotFound)
}

// DeleteEvent deletes an event with the specified ID
func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, event := range events {
		if strconv.Itoa(event.ID) == params["id"] {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "Event with ID %s is deleted successfully", params["id"])
			return
		}
	}
	http.Error(w, "Event not found", http.StatusNotFound)
}

// Register creates a new user account
func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, ok := users[user.Name]; ok {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}
	user.ID = len(users) + 1
	users[strconv.Itoa(user.ID)] = user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// FindUserByEmail finds a user by email
func FindUserByEmail(email string) (User, bool) {
	for _, user := range users {
		if user.Email == email {
			return user, true
		}
	}
	return User{}, false
}

// Login authenticates a user and returns a JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := FindUserByEmail(credentials.Email)

	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if user.Password != credentials.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{Token: tokenString})
}

func main() {
	router := mux.NewRouter()
	// Event endpoints
	router.HandleFunc("/events", AuthMiddleware(CreateEvent)).Methods("POST")
	router.HandleFunc("/events/{id}", GetEvent).Methods("GET")
	router.HandleFunc("/events", GetAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", AuthMiddleware(UpdateEvent)).Methods("PUT")
	router.HandleFunc("/events/{id}", AuthMiddleware(DeleteEvent)).Methods("DELETE")

	// User endpoints
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
