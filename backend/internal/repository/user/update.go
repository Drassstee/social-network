package user

import "social-network/internal/models/user"

func (r *UserRepo) UpdateUser(u *user.User) error {
	query := `UPDATE users
			SET first_name = ?, last_name = ?, date_of_birth = ?, nickname = ?, about_me = ?, profile_type = ?
			WHERE id = ?`

	_, err := r.db.Exec(query,
		u.FirstName,
		u.LastName,
		u.DOB,
		u.Nickname,
		u.AboutMe,
		u.ProfileType,
		u.ID,
	)
	return err
}

func (r *UserRepo) UpdateAvatar(id int64, url string) error {
	query := `UPDATE users
			SET avatar_url = ?
			WHERE id = ?`

	_, err := r.db.Exec(query, url, id)
	return err
}
