package user

import (
	"social-network/internal/models/follow"
	"social-network/internal/models/post"
	"social-network/internal/models/session"
	"social-network/internal/models/user"
)

type UserService struct {
	users    UserRepo
	sessions SessionRepo
	follows  FollowRepo
	posts    PostRepo
}

type UserRepo interface {
	CreateUser(*user.User) (int64, error)
	UpdateUser(*user.User) error

	GetByEmail(string) (*user.User, error)
	GetByID(int64) (*user.User, error)
	GetProfileType(int64) (string, error)

	IsPrivate(int64) (bool, error)
	EmailExists(string, int64) (bool, error)
	UserExists(int64) (bool, error)
}

type SessionRepo interface {
	CreateSession(*session.Session) error
	DeleteSession(string) error
	DeleteAllSessions(int64) error

	GetUserID(string) (int64, error)
	GetUUID(int64) (string, error)
}

type FollowRepo interface {
	CreateFollow(follow.Follow) error
	DeleteFollow(follow.Follow) error
	UpdateFollow(follow.Follow) error

	GetFollowers(int64, string) ([]user.UserData, error)
	GetFollowing(int64, string) ([]user.UserData, error)

	IsFollower(int64, int64) (bool, error)
	FollowExists(int64, int64, string) (bool, error)
}

type PostRepo interface {
	GetPosts(int64) ([]post.Post, error)
}

func NewUserService(ur UserRepo, sr SessionRepo, fr FollowRepo, pr PostRepo) *UserService {
	return &UserService{
		ur, sr, fr, pr,
	}
}
