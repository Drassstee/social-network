package session

func (r *SessionRepo) DeleteSession(uuid string) error {
	query := `DELETE FROM sessions WHERE uuid = ?`

	_, err := r.db.Exec(query, uuid)
	return err
}

// --------------------------------------------------------------------|

func (r *SessionRepo) DeleteAllSessions(id int64) error {
	query := `DELETE FROM sessions WHERE user_id = ?`

	_, err := r.db.Exec(query, id)
	return err
}
