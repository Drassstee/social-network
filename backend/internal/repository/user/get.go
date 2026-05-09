package user

import (
	"context"
	"database/sql"
	"fmt"

	"social-network/internal/models"
)

// --------------------------------------------------------------------|

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, email, first_name, last_name, date_of_birth, 
			COALESCE(nickname, ''), COALESCE(about_me, ''), profile_type, COALESCE(avatar_url, '')
			FROM users
			WHERE id = ?`

	var u models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.DOB, &u.Nickname, &u.AboutMe, &u.ProfileType, &u.AvatarURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: user not found", models.ErrNotFound)
		}
		return nil, err
	}
	return &u, nil
}

// --------------------------------------------------------------------|

func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, password, email, 
			COALESCE(nickname, ''), COALESCE(avatar_url, ''), date_of_birth, COALESCE(about_me, ''), profile_type
			FROM users
			WHERE email = ?`

	var u models.User
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Password, &u.Email, &u.Nickname, &u.AvatarURL, &u.DOB, &u.AboutMe, &u.ProfileType)
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
	query := `SELECT COALESCE(avatar_url, '') FROM users WHERE id = ?`

	var url string
	err := r.db.QueryRow(query, id).Scan(&url)
	return url, err
}

// --------------------------------------------------------------------|

func (r *UserRepo) GetByIDs(ids []int64) ([]models.User, error) {
	if len(ids) == 0 {
		return []models.User{}, nil
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

	query := `SELECT id, email, first_name, last_name, date_of_birth, 
			COALESCE(avatar_url, ''), COALESCE(nickname, ''), COALESCE(about_me, ''), profile_type 
			FROM users WHERE id IN (` + placeholders + `)`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
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
