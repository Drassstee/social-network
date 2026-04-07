package user

import (
	"fmt"

	"social-network/internal/models"
	"social-network/internal/models/profile"
	"social-network/internal/models/user"
)

func (us *UserService) GetProfile(userID, targetID int64) (*profile.Profile, error) {
	if targetID < 1 {
		return nil, fmt.Errorf("get profile: %w: incorrect user id", models.ErrInvalidData)
	}

	if userID != targetID {
		isPrivate, err := us.users.IsPrivate(targetID)
		if err != nil {
			return nil, fmt.Errorf("get profile: %w", err)
		}

		if isPrivate {
			isFollower, err := us.follows.IsFollower(userID, targetID)
			if err != nil {
				return nil, fmt.Errorf("get profile: %w", err)
			}

			if !isFollower {
				data, err := us.users.GetByID(targetID)
				if err != nil {
					return nil, fmt.Errorf("get profile: %w", err)
				}

				var u = user.User{
					ID:          data.ID,
					FirstName:   data.FirstName,
					LastName:    data.LastName,
					ProfileType: data.ProfileType,
				}

				return &profile.Profile{
					User: &u,
				}, nil
			}
		}
	}

	u, err := us.users.GetByID(targetID)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	posts, err := us.posts.GetPosts(targetID)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	followers, err := us.follows.GetFollowers(targetID, "accept")
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	following, err := us.follows.GetFollowing(targetID, "accept")
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	return &profile.Profile{
		User:      u,
		Posts:     posts,
		Followers: followers,
		Following: following,
	}, nil
}

// --------------------------------------------------------------------|

func (us *UserService) UpdateProfile(u *user.User) error {
	err := u.ValidateData()
	if err != nil {
		return fmt.Errorf("update profile: %w: %w", models.ErrInvalidData, err)

	}

	if u.ProfileType != "public" && u.ProfileType != "private" {
		return fmt.Errorf("update profile: %w: incorrect profile type", models.ErrInvalidData)
	}

	exists, err := us.users.EmailExists(u.Email, u.ID)
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}
	if exists {
		return fmt.Errorf("get profile: %w: email already exists", models.ErrConflict)
	}

	if err := us.users.UpdateUser(u); err != nil {
		return fmt.Errorf("update profile: %w", err)
	}
	return nil
}

// --------------------------------------------------------------------|

func (us *UserService) GetUserID(uuid string) (int64, error) {
	id, err := us.sessions.GetUserID(uuid)
	if err != nil {
		return 0, fmt.Errorf("get user id: %w", err)
	}

	return id, nil
}
