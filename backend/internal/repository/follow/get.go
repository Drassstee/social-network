package follow

import (
	"log"

	"social-network/internal/models/user"
)

func (r *FollowRepo) GetFollowers(id int64, status string) ([]user.UserData, error) {
	query := `SELECT u.id, u.first_name, u.last_name
			FROM follows AS f 
			LEFT JOIN users AS u ON u.id = f.follower_id
			WHERE f.status = ? AND f.following_id = ?`

	rows, err := r.db.Query(query, status, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.UserData
	for rows.Next() {
		var u user.UserData
		err = rows.Scan(&u.ID, &u.FirstName, &u.LastName)
		if err != nil {
			log.Printf("scan user: %v", err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

// --------------------------------------------------------------------|

func (r *FollowRepo) GetFollowing(id int64, status string) ([]user.UserData, error) {
	query := `SELECT u.id, u.first_name, u.last_name
			FROM follows AS f 
			LEFT JOIN users AS u ON u.id = f.following_id
			WHERE f.status = ? AND f.follower_id = ?`

	rows, err := r.db.Query(query, status, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.UserData
	for rows.Next() {
		var u user.UserData
		err = rows.Scan(&u.ID, &u.FirstName, &u.LastName)
		if err != nil {
			log.Printf("scan user: %v", err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

// --------------------------------------------------------------------|
