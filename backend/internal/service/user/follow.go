package user

import (
	"fmt"

	"social-network/internal/models"
	"social-network/internal/models/follow"
)

func (us *UserService) Follow(f follow.Follow) (string, error) {
	if f.FollowerID < 1 || f.FollowingID < 1 {
		return "", fmt.Errorf("follow: %w: incorrect user id", models.ErrInvalidData)
	}

	if f.FollowerID == f.FollowingID {
		return "", fmt.Errorf("follow: %w: self-following is not allowed", models.ErrInvalidData)
	}

	exists, err := us.users.UserExists(f.FollowingID)
	if err != nil {
		return "", fmt.Errorf("follow: %w", err)
	}

	if !exists {
		return "", fmt.Errorf("follow: %w: incorrect user id", models.ErrInvalidData)
	}

	ptype, err := us.users.GetProfileType(f.FollowingID)
	if err != nil {
		return "", fmt.Errorf("follow: %w", err)
	}

	if ptype == "public" {
		f.Status = "accept"
	} else {
		f.Status = "pending"
	}

	err = us.follows.CreateFollow(f)
	if err != nil {
		return "", fmt.Errorf("follow: %w", err)
	}

	return f.Status, nil
}

// --------------------------------------------------------------------|

func (us *UserService) Unfollow(f follow.Follow) error {
	if f.FollowerID < 1 || f.FollowingID < 1 {
		return fmt.Errorf("follow: %w: incorrect user id", models.ErrInvalidData)
	}

	if f.FollowerID == f.FollowingID {
		return fmt.Errorf("follow: %w: self-unfollowing is not allowed", models.ErrInvalidData)
	}

	exists, err := us.follows.FollowExists(f.FollowerID, f.FollowingID, "accept")
	if err != nil {
		return fmt.Errorf("unfollow: %w", err)
	}

	if !exists {
		return fmt.Errorf("unfollow: %w: not following this user", models.ErrInvalidData)
	}

	err = us.follows.DeleteFollow(f)
	if err != nil {
		return fmt.Errorf("unfollow: %w", err)
	}

	return nil
}
