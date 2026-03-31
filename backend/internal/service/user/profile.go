package user

import (
	"fmt"

	"social-network/internal/models/profile"
	"social-network/internal/models/user"
)

func (r *UserService) GetProfile(id int64) (*profile.Profile, error) {
	if id < 1 {
		return nil, user.ErrIncorrectID
	}

	u, err := r.users.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err) // err: user not found || db.err
	}

	posts, err := r.posts.GetPosts(id)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err) // err: db.err
	}

	followers, err := r.follows.GetFollowers(id)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err) // err: db.err
	}

	following, err := r.follows.GetFollowing(id)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err) // err: db.err
	}

	return &profile.Profile{
		User:      u,
		Posts:     posts,
		Followers: followers,
		Following: following,
	}, nil
}

// --------------------------------------------------------------------|

func (r *UserService) UpdateProfile(u *user.User) error {
	err := u.ValidateData()
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}

	if err := r.users.UpdateUser(u); err != nil {
		return fmt.Errorf("update profile: %w", err) // err: db.err
	}
	return nil
}

// --------------------------------------------------------------------|

func (r *UserService) GetUserID(uuid string) (int64, error) {
	id, err := r.sessions.GetUserID(uuid)
	if err != nil {
		return 0, fmt.Errorf("get user id: %w", err) // err: user not found || db.err
	}

	return id, nil
}
