package follow

import "social-network/internal/models/follow"

func (r *FollowRepo) DeleteFollow(f follow.Follow) error {
	query := `DELETE FROM follows WHERE follower_id = ? AND following_id = ?`

	_, err := r.db.Exec(query, f.FollowerID, f.FollowingID)
	return err
}
