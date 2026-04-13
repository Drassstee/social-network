package user

import (
	"database/sql"
	"fmt"

	"social-network/internal/models"
	"social-network/internal/models/user"
)

func (r *UserRepo) GetByID(id int64) (*user.User, error) {
	query := `SELECT id, email, first_name, last_name, date_of_birth, nickname, about_me, profile_type
			FROM users
			WHERE id = ?`

	var u user.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.DOB, &u.Nickname, &u.AboutMe, &u.ProfileType)
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
	return ptype, err
}

// --------------------------------------------------------------------|

func (r *UserRepo) GetAvatarURL(id int64) (string, error) {
	query := `SELECT avatar_url FROM users WHERE id = ?`

	var url string
	err := r.db.QueryRow(query, id).Scan(&url)
	return url, err
}

// --------------------------------------------------------------------|

func (r *UserRepo) GetByIDs(ids []int64) ([]user.User, error) {
	if len(ids) == 0 {
		return []user.User{}, nil
	}

	placeholders := ""
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = id
	}

	query := `SELECT id, email, first_name, last_name, date_of_birth, avatar_url, nickname, about_me, profile_type FROM users WHERE id IN (` + placeholders + `)`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.DOB, &u.AvatarURL, &u.Nickname, &u.AboutMe, &u.ProfileType); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
