package user

import (
	"context"
	"fmt"

	"social-network/internal/models"
	"social-network/internal/utils"
)

//--------------------------------------------------------------------------------------|

func (s *UserService) GetProfile(targetID, userID int64) (*models.Profile, error) {
	if targetID < 1 {
		return nil, fmt.Errorf("%w: target user id is empty", models.ErrInvalidData)
	}
	ctx := context.Background()

	u, err := s.users.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	p := &models.Profile{
		User: models.UserData{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Nickname:  u.Nickname,
			AvatarURL: utils.FormatAvatarURL(u.AvatarURL),
		},
		Privacy: u.ProfileType,
		DOB:     u.DOB,
		AboutMe: u.AboutMe,
		Email:   u.Email,
	}

	if userID == targetID {
		p.IsMe = true
	} else {
		// Check for any existing follow relationship
		status, err := s.follows.GetFollowStatus(userID, targetID)
		if err == nil && status != "" {
			if status == "accept" {
				p.FollowStatus = "following"
			} else if status == "pending" {
				p.FollowStatus = "pending"
			}
		}

		if u.ProfileType == "private" && p.FollowStatus != "following" {
			return p, models.ErrUserPrivate
		}
	}

	// Fetch Followers
	followers, err := s.follows.GetFollowers(targetID, "accept")
	if err == nil {
		p.Followers = followers
	}

	// Fetch Following
	following, err := s.follows.GetFollowing(targetID, "accept")
	if err == nil {
		p.Following = following
	}

	// Fetch Posts
	posts, err := s.posts.GetPosts(targetID, userID)
	if err == nil {
		for i := range posts {
			if posts[i].Author != nil {
				posts[i].Author.AvatarURL = utils.FormatAvatarURL(posts[i].Author.AvatarURL)
			}
		}
		p.Posts = posts
	}

	return p, nil
}

//--------------------------------------------------------------------------------------|

func (s *UserService) GetMe(id int64) (*models.UserData, error) {
	if id < 1 {
		return nil, fmt.Errorf("%w: incorrect user id", models.ErrInvalidData)
	}

	u, err := s.users.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &models.UserData{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		AvatarURL: utils.FormatAvatarURL(u.AvatarURL),
	}, nil
}

//--------------------------------------------------------------------------------------|

func (s *UserService) UpdateProfile(u *models.User) error {
	if err := u.ValidateData(); err != nil {
		return fmt.Errorf("%w: %v", models.ErrInvalidData, err)
	}

	exists, err := s.users.EmailExists(u.Email, u.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: email already exists", models.ErrConflict)
	}

	return s.users.UpdateUser(u)
}

//--------------------------------------------------------------------------------------|

func (s *UserService) GetUserID(uuid string) (int64, error) {
	return s.sessions.GetUserID(uuid)
}
