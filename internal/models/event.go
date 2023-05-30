package models

type Event struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Date    string `json:"date"`
	Details string `json:"details"`
}
