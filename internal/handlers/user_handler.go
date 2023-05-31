package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/TheVitik/sport-together-go/internal/models"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

var secretKey = []byte("secret")

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	fmt.Println("START")
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.repository.FindUserByEmail(credentials.Email)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if user.Password != credentials.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := newToken.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	var token models.Token
	token.Token = tokenString

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tempUser, err := h.repository.FindUserByEmail(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tempUser != nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	user, err = h.repository.SaveUser(user)
	if err != nil {
		http.Error(w, "Save user error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Panic(err)
		return
	}
}
