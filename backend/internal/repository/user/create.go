package user

import (
	"errors"
	"fmt"

	"social-network/internal/models"
	"social-network/internal/models/user"

	"github.com/mattn/go-sqlite3"
)

func (r *UserRepo) CreateUser(u *user.User) (int64, error) {
	query := `INSERT INTO 
			users (email, first_name, last_name, password, date_of_birth, avatar_url, nickname, about_me) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Password,
		u.DOB,
		u.AvatarURL,
		u.Nickname,
		u.AboutMe,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%w: email already exists", models.ErrConflict)
		}
		return 0, err
	}
	return res.LastInsertId()
}
