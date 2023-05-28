package app

import (
	"database/sql"
	"encoding/json"
	ihttp "github.com/TheVitik/sport-together-go/internal/http"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// secretKey is the secret key used to sign JWT tokens
var secretKey = []byte("secret")

// Login authenticates a user and returns a JWT token
func Login(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := findUserByEmail(credentials.Email, db)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
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
		json.NewEncoder(w).Encode(ihttp.Token{Token: tokenString})
	}
	return handler
}

// Register creates a new user account
func Register(db *sql.DB) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var user ihttp.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tempUser, err := findUserByEmail(user.Email, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if tempUser != nil {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
		row := db.QueryRow(query, user.Name, user.Email, user.Password)

		err = row.Scan(&user.ID)
		if err != nil {
			http.Error(w, "Save user error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
	return handler
}

// FindUserByEmail finds a user by email
func findUserByEmail(email string, db *sql.DB) (*ihttp.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	row := db.QueryRow(query, email)

	user := &ihttp.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}

// findToken finds a token in db
func findToken(token string, db *sql.DB) (*ihttp.User, error) {
	query := "SELECT * FROM tokens WHERE token = $1"
	row := db.QueryRow(query, token)

	user := &ihttp.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}
