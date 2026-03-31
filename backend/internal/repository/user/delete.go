package user

func (r *UserRepo) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	return err
}
