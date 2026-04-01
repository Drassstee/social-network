// Package sessions provides database-backed session management
// with automatic expiration and cleanup.
package sessions

import (
	"context"
	"database/sql"
	"errors"
	"forum/internal/models"
	"time"

	"github.com/google/uuid"
)

//--------------------------------------------------------------------------------------|

// DefaultSessionTTL defines the default duration (2 hours) before a session expires.
const DefaultSessionTTL = 2 * time.Hour

//--------------------------------------------------------------------------------------|

// DBSessionStorage implements session persistence using a SQL database.
type DBSessionStorage struct {
	db *sql.DB
}

// SessionManager handles the high-level logic for creating, retrieving, and
// deleting user sessions with a configurable Time-To-Live (TTL).
type SessionManager struct {
	storage *DBSessionStorage
	ttl     time.Duration
}

// NewSessionManager creates a new SessionManager with the provided database
// connection and session duration. If ttl is 0, DefaultSessionTTL is used.
func NewSessionManager(db *sql.DB, ttl time.Duration) *SessionManager {
	if ttl == 0 {
		ttl = DefaultSessionTTL
	}
	return &SessionManager{
		storage: &DBSessionStorage{db: db},
		ttl:     ttl,
	}
}

//--------------------------------------------------------------------------------------|

// CreateSession initiates a new session for the given user, invalidating any existing sessions.
func (sm *SessionManager) CreateSession(ctx context.Context, userID int) (*models.Session, error) {
	return sm.storage.CreateSession(ctx, userID, sm.ttl)
}

// GetSession retrieves an active session by its ID. It returns an error if the
// session is expired or does not exist.
func (sm *SessionManager) GetSession(ctx context.Context, id string) (*models.Session, error) {
	return sm.storage.GetSession(ctx, id)
}

// DeleteSession invalidates/removes the session identified by the given ID.
func (sm *SessionManager) DeleteSession(ctx context.Context, id string) error {
	return sm.storage.DeleteSession(ctx, id)
}

//--------------------------------------------------------------------------------------|

// CreateSession generates a new UUID-based session, invalidates old sessions for the user,
// and saves the new session to the database.
func (s *DBSessionStorage) CreateSession(ctx context.Context, userID int, ttl time.Duration) (*models.Session, error) {
	id := uuid.New().String()

	now := time.Now().UTC().Truncate(time.Second)
	expiresAt := now.Add(ttl).Truncate(time.Second)

	_, err := s.db.ExecContext(ctx,
		`DELETE FROM sessions WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO sessions (id, user_id, created_at, expires_at) 
         VALUES (?, ?, ?, ?)`,
		id, userID, now, expiresAt)
	if err != nil {
		return nil, err
	}

	return &models.Session{
		ID:        id,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}, nil
}

//--------------------------------------------------------------------------------------|

// GetSession looks up a session in the database. It enforces expiration checks using the current time.
func (s *DBSessionStorage) GetSession(ctx context.Context, id string) (*models.Session, error) {
	var session models.Session
	err := s.db.QueryRowContext(ctx,
		`SELECT id, user_id, created_at, expires_at 
         FROM sessions WHERE id = ? AND expires_at > ?`,
		id, time.Now().UTC()).Scan(
		&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

//--------------------------------------------------------------------------------------|

// DeleteSession removes the session record from the database.
func (s *DBSessionStorage) DeleteSession(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE id = ?`, id)
	return err
}
