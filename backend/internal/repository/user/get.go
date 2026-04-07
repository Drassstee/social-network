package user

import (
	"database/sql"
	"fmt"

	"social-network/internal/models"
	"social-network/internal/models/user"
)

func (r *UserRepo) GetByID(id int64) (*user.User, error) {
	query := `SELECT id, email, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type
			FROM users
			WHERE id = ?`

	var u user.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.DOB, &u.Avatar, &u.Nickname, &u.AboutMe, &u.ProfileType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: user not found", models.ErrNotFound)
		}
		return nil, err
	}
	return &u, nil
}

// --------------------------------------------------------------------|

func (r *UserRepo) GetByEmail(email string) (*user.User, error) {
	query := `SELECT id, first_name, last_name, password
			FROM users
			WHERE email = ?`

	var u user.User
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: invalid email or password", models.ErrInvalidData)
		}
		return nil, err
	}
	return &u, nil
}

// --------------------------------------------------------------------|

func (r *UserRepo) GetProfileType(id int64) (string, error) {
	query := `SELECT profile_type FROM users WHERE id = ?`

	var ptype string
	err := r.db.QueryRow(query, id).Scan(&ptype)
	if err != nil {
		return "", err
	}
	return ptype, nil
}
