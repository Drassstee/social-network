package follow

import (
	"errors"
	"fmt"

	"social-network/internal/models"
	"social-network/internal/models/follow"

	"github.com/mattn/go-sqlite3"
)

func (r *FollowRepo) CreateFollow(f follow.Follow) error {
	query := `INSERT INTO follows (follower_id, following_id, status) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, f.FollowerID, f.FollowingID, f.Status)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			return fmt.Errorf("%w: follow already exists", models.ErrConflict)
		}
		return err
	}
	return nil
}
