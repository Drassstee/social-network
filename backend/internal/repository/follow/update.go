package follow

import "social-network/internal/models"

func (r *FollowRepo) UpdateFollow(f models.Follow) error {
	query := `UPDATE follows SET status = ? WHERE follower_id = ? AND following_id = ?`

	_, err := r.db.Exec(query, f.Status, f.FollowerID, f.FollowingID)
	return err
}
