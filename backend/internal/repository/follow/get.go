package follow

import (
	"database/sql"
	"log"

	"social-network/internal/models"
)

func (r *FollowRepo) GetFollowers(id int64, status string) ([]models.UserData, error) {
	query := `SELECT u.id, u.first_name, u.last_name, u.avatar_url, u.nickname
			FROM follows AS f 
			LEFT JOIN users AS u ON u.id = f.follower_id
			WHERE f.status = ? AND f.following_id = ?`

	rows, err := r.db.Query(query, status, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserData
	for rows.Next() {
		var u models.UserData
		err = rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.AvatarURL, &u.Nickname)
		if err != nil {
			log.Printf("scan user: %v", err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

// --------------------------------------------------------------------|

func (r *FollowRepo) GetFollowing(id int64, status string) ([]models.UserData, error) {
	query := `SELECT u.id, u.first_name, u.last_name, u.avatar_url, u.nickname
			FROM follows AS f 
			LEFT JOIN users AS u ON u.id = f.following_id
			WHERE f.status = ? AND f.follower_id = ?`

	rows, err := r.db.Query(query, status, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserData
	for rows.Next() {
		var u models.UserData
		err = rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.AvatarURL, &u.Nickname)
		if err != nil {
			log.Printf("scan user: %v", err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *FollowRepo) GetStatus(followerID, followingID int64) (string, error) {
	var status string
	err := r.db.QueryRow(`SELECT status FROM follows WHERE follower_id = ? AND following_id = ?`, followerID, followingID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return status, nil
}

// --------------------------------------------------------------------|
