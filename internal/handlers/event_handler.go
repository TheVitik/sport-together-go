package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/TheVitik/sport-together-go/internal/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	event, err = h.repository.SaveEvent(event)
	if err != nil {
		http.Error(w, "Save event error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		return
	}
}

// GetEvent returns an event with the specified ID
func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	event, err := h.repository.GetEvent(params["id"])
	if err != nil {
		http.Error(w, "Save event error", http.StatusInternalServerError)
		return
	}
	if event == nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		return
	}

}

func (h *Handler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.repository.GetEvents()
	if err != nil {
		http.Error(w, "Get events error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		return
	}
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedEvent *models.Event
	err := json.NewDecoder(r.Body).Decode(&updatedEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}
	updatedEvent.ID = id
	event, err := h.repository.UpdateEvent(updatedEvent)
	if err != nil {
		http.Error(w, "Update event error", http.StatusInternalServerError)
		return
	}
	if event == nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		return
	}
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	event, err := h.repository.DeleteEvent(params["id"])
	if err != nil {
		http.Error(w, "Delete event error", http.StatusInternalServerError)
		return
	}
	if event == false {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	_, err = fmt.Fprintf(w, "Event with ID %s is deleted successfully", params["id"])
	if err != nil {
		return
	}
}
