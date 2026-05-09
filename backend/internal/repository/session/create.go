package session

import "social-network/internal/models"

//--------------------------------------------------------------------------------------|

func (r *SessionRepo) CreateSession(s *models.Session) error {
	query := `INSERT INTO sessions (uuid, user_id, expires_at) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, s.ID, s.UserID, s.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}
