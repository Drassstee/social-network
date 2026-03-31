package session

import (
	"database/sql"
	"errors"

	"social-network/internal/models/session"
)

func (r *SessionRepo) GetUserID(uuid string) (int64, error) {
	query := `SELECT user_id
			FROM sessions
			WHERE uuid = ? AND expires_at > DATETIME('now')`

	var id int64
	err := r.db.QueryRow(query, uuid).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, session.ErrNotFound
		}
		return 0, err
	}
	return id, nil
}
