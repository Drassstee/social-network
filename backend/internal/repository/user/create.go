package user

import (
	"social-network/internal/models"
)

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, email, first_name, last_name, password) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, user.ID, user.Email, user.FirstName, user.LastName, user.Password)
	return err
}
