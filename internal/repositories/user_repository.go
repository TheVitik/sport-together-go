package repositories

import (
	"database/sql"
	"github.com/TheVitik/sport-together-go/internal/models"
)

func (r *Repository) SaveUser(user models.User) (models.User, error) {
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	row := r.db.QueryRow(query, user.Name, user.Email, user.Password)

	var err = row.Scan(&user.ID)
	return user, err
}

func (r *Repository) FindUserByEmail(email string) (*models.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	row := r.db.QueryRow(query, email)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
