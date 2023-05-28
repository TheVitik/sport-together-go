package app

import (
	"encoding/json"
	"fmt"
	ihttp "github.com/TheVitik/sport-together-go/internal/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// events is a slice of Event structs
var events []ihttp.Event

// GetAllEvents returns all events in the events slice
func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

// CreateEvent creates a new event and adds it to the events slice
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event ihttp.Event
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

// UpdateEvent updates an event with the specified ID
func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, event := range events {
		if strconv.Itoa(event.ID) == params["id"] {
			events = append(events[:i], events[i+1:]...)
			var updatedEvent ihttp.Event
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
