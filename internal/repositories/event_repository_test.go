package repositories

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TheVitik/sport-together-go/internal/database"
	"github.com/TheVitik/sport-together-go/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveEvent(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repository := NewRepository(&database.Connection{DB: db})

	dbMock.ExpectExec("INSERT INTO events").
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repository.SaveEvent(&models.Event{Name: "Test", Date: "31-05-23", Details: "description"})
	assert.Nil(t, err)

	err = dbMock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "Test", result.Name)
}

func TestGetEvent(t *testing.T) {
	db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer db.Close()

	repository := NewRepository(&database.Connection{DB: db})

	columns := []string{"id", "name", "date", "details"}
	event := models.Event{ID: 1, Name: "Test", Date: "31-05-23", Details: "description"}
	dbMock.ExpectQuery("SELECT * FROM events WHERE id = $1 LIMIT 1").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(event.ID, event.Name, event.Date, event.Details))

	result, err := repository.GetEvent("1")
	assert.Nil(t, err)

	err = dbMock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, event.ID, result.ID)
	assert.Equal(t, event.Name, result.Name)
	assert.Equal(t, event.Date, result.Date)
	assert.Equal(t, event.Details, result.Details)
}

func TestUpdateEvent(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repository := NewRepository(&database.Connection{DB: db})

	dbMock.ExpectExec("INSERT INTO events").
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repository.SaveEvent(&models.Event{Name: "Test", Date: "31-05-23", Details: "description"})
	assert.Nil(t, err)

	dbMock.ExpectExec("UPDATE events").
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err = repository.UpdateEvent(&models.Event{ID: 1, Name: "Test update", Date: "31-05-23", Details: "description"})
	assert.Nil(t, err)

	err = dbMock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "Test update", result.Name)
}

func TestDeleteEvent(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repository := NewRepository(&database.Connection{DB: db})

	dbMock.ExpectExec("INSERT INTO events").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repository.SaveEvent(&models.Event{Name: "Test", Date: "31-05-23", Details: "description"})
	assert.Nil(t, err)

	dbMock.ExpectExec("DELETE FROM events").
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := repository.DeleteEvent("1")
	assert.Nil(t, err)

	err = dbMock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, true, res)
}
