package user

import (
	"fmt"

	"social-network/internal/models"
	"social-network/internal/models/follow"
	"social-network/internal/models/user"
)

func (us *UserService) GetNotification(id int64) ([]user.UserData, error) {
	if id < 1 {
		return nil, fmt.Errorf("notification: %w: incorrect user id", models.ErrInvalidData)
	}

	users, err := us.follows.GetFollowers(id, "pending")
	if err != nil {
		return nil, fmt.Errorf("notification: %w", err)
	}

	return users, nil
}

// --------------------------------------------------------------------|

func (us *UserService) RespondToFollowRequest(f follow.Follow) error {
	if f.FollowerID < 1 || f.FollowingID < 1 {
		return fmt.Errorf("notification: %w: incorrect user id", models.ErrInvalidData)
	}

	if f.FollowerID == f.FollowingID {
		return fmt.Errorf("notification: %w: same id", models.ErrInvalidData)
	}

	exists, err := us.follows.FollowExists(f.FollowerID, f.FollowingID, "pending")
	if err != nil {
		return fmt.Errorf("notification: %w", err)
	}

	if !exists {
		return fmt.Errorf("notification: %w: follow request not found", models.ErrNotFound)
	}

	switch f.Status {
	case "accept":
		err := us.follows.UpdateFollow(f)
		if err != nil {
			return fmt.Errorf("notification: %w", err)
		}
	case "decline":
		err := us.follows.DeleteFollow(f)
		if err != nil {
			return fmt.Errorf("notification: %w", err)
		}
	default:
		return fmt.Errorf("notification: %w: incorrect status", models.ErrInvalidData)
	}

	return nil
}
