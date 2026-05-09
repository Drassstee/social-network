package follow

func (r *FollowRepo) IsFollower(userID, targetID int64) (bool, error) {
	query := `SELECT EXISTS(
				SELECT 1 
				FROM follows 
				WHERE follower_id = ? AND following_id = ? AND status = 'accept')`


	var exists int
	err := r.db.QueryRow(query, userID, targetID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// --------------------------------------------------------------------|

func (r *FollowRepo) FollowExists(followerID, followingID int64, status string) (bool, error) {
	query := `SELECT EXISTS(
				SELECT 1
				FROM follows
				WHERE follower_id = ? AND following_id = ? AND status = ?)`

	var exists int
	err := r.db.QueryRow(query, followerID, followingID, status).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// --------------------------------------------------------------------|

func (r *FollowRepo) GetFollowStatus(followerID, followingID int64) (string, error) {
	query := `SELECT status FROM follows WHERE follower_id = ? AND following_id = ?`

	var status string
	err := r.db.QueryRow(query, followerID, followingID).Scan(&status)
	if err != nil {
		return "", nil // Return empty if not found, no error needed for this specific logic
	}
	return status, nil
}
