package user

func (r *UserRepo) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	return err
}
