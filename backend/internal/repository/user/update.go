package user

import "social-network/internal/models/user"

func (r *UserRepo) UpdateUser(u *user.User) error {
	query := `UPDATE users
			SET first_name = ?, last_name = ?, date_of_birth = ?, avatar = ?, nickname = ?, about_me = ?, profile_type = ?
			WHERE id = ?`

	_, err := r.db.Exec(query,
		u.FirstName,
		u.LastName,
		u.DOB,
		u.Avatar,
		u.Nickname,
		u.AboutMe,
		u.ProfileType,
		u.ID,
	)
	return err
}
