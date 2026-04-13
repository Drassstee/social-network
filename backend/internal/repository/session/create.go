package session

import "social-network/internal/models/session"

func (r *SessionRepo) CreateSession(s *session.Session) error {
	query := `INSERT INTO sessions (uuid, user_id, expires_at) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, s.UUID, s.UserID, s.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}
