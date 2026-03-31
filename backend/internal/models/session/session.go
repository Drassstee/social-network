package session

import "time"

type Session struct {
	UUID      string
	UserID    int64
	ExpiresAt time.Time
}
