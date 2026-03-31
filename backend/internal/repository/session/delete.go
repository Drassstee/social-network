package session

func (r *SessionRepo) Delete(uuid string) error {
	query := `DELETE FROM sessions WHERE uuid = ?`

	_, err := r.db.Exec(query, uuid)
	return err
}
