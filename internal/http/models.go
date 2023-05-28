package http

// User is a struct representing a user
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Event is a struct representing an event
type Event struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Date    string `json:"date"`
	Details string `json:"details"`
}

// Token is a struct representing an auth token
type Token struct {
	Token string `json:"token"`
}
