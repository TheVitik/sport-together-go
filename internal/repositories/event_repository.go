package repositories

import (
	"database/sql"
	"github.com/TheVitik/sport-together-go/internal/models"
	"log"
)

func (r *Repository) GetEvents() ([]models.Event, error) {
	query := "SELECT * FROM events"
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	var events []models.Event
	for rows.Next() {
		var event models.Event

		// Scan the row values into the Event struct fields
		err := rows.Scan(&event.ID, &event.Name, &event.Date, &event.Details)
		if err != nil {
			log.Fatal(err)
		}

		// Append the event to the events slice
		events = append(events, event)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return events, nil
}

func (r *Repository) SaveEvent(event models.Event) (models.Event, error) {
	query := "INSERT INTO events (name, date, details) VALUES ($1, $2, $3) RETURNING id"
	row := r.db.QueryRow(query, event.Name, event.Date, event.Details)

	var err = row.Scan(&event.ID)
	return event, err
}

func (r *Repository) GetEvent(id string) (*models.Event, error) {
	query := "SELECT * FROM events WHERE id = $1 LIMIT 1"
	row := r.db.QueryRow(query, id)

	event := &models.Event{}
	err := row.Scan(&event.ID, &event.Name, &event.Date, &event.Details)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return event, nil
}

func (r *Repository) UpdateEvent(event *models.Event) (*models.Event, error) {
	query := "UPDATE events SET name=$1, date=$2, details=$3 WHERE id=$4"
	result, err := r.db.Exec(query, event.Name, event.Date, event.Details, event.ID)
	if err != nil {
		return nil, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	return event, err
}

func (r *Repository) DeleteEvent(id string) (bool, error) {
	query := "DELETE FROM events WHERE id=$1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}

	return true, err
}
