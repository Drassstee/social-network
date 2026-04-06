package user

import (
	"context"
	"database/sql"
	"social-network/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	query := `SELECT id, email, username, first_name, last_name, password FROM users WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query := `SELECT id, email, username, first_name, last_name, password FROM users WHERE email = ?`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *UserRepository) GetByIDs(ctx context.Context, ids []int) ([]models.User, error) {
	if len(ids) == 0 {
		return []models.User{}, nil
	}

	query := `SELECT id, email, username, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type FROM users WHERE id IN (`
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		query += "?"
		if i < len(ids)-1 {
			query += ","
		}
		args[i] = id
	}
	query += ")"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName, &u.BirthDay, &u.Avatar, &u.Nickname, &u.Bio, &u.ProfileType); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
