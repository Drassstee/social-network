package session

import (
	"database/sql"
	"errors"
	"fmt"

	"social-network/internal/models"
)

func (r *SessionRepo) GetUserID(uuid string) (int64, error) {
	query := `SELECT user_id
			FROM sessions
			WHERE uuid = ? AND expires_at > DATETIME('now')`

	var id int64
	err := r.db.QueryRow(query, uuid).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no such session: %w", models.ErrNotFound)
		}
		return 0, err
	}
	return id, nil
}

// --------------------------------------------------------------------|

func (r *SessionRepo) GetUUID(id int64) (string, error) {
	query := `SELECT uuid
			FROM sessions
			WHERE user_id = ?`

	var uuid string
	err := r.db.QueryRow(query, id).Scan(&uuid)
	if err != nil {
		return "", err
	}
	return uuid, nil
}
