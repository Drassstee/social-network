package follow

import "social-network/internal/models/follow"

func GetFollowers(id int64) ([]follow.Follow, error) {
	query := `SELECT `
}
