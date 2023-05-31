package repositories

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TheVitik/sport-together-go/internal/database"
	"github.com/TheVitik/sport-together-go/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveUser(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repository := NewRepository(&database.Connection{DB: db})

	dbMock.ExpectExec("INSERT INTO users").
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repository.SaveUser(&models.User{Name: "Test", Email: "test@mail.com", Password: "11111111"})
	assert.Nil(t, err)

	err = dbMock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "Test", result.Name)
}

func TestFindUserByEmail(t *testing.T) {
	db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer db.Close()

	repository := NewRepository(&database.Connection{DB: db})

	columns := []string{"id", "name", "email", "password"}
	user := models.User{ID: 1, Name: "Test", Email: "test@mail.com", Password: "11111111"}
	dbMock.ExpectQuery("SELECT * FROM users WHERE email = $1").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(user.ID, user.Name, user.Email, user.Password))

	result, err := repository.FindUserByEmail("test@mail.com")
	assert.Nil(t, err)

	err = dbMock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.Password, result.Password)
}
