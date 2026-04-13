package models

import (
	"social-network/internal/models/follow"
	"social-network/internal/models/session"
	"social-network/internal/models/user"
)

type User = user.User

type UserService interface {
	GetByID(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByIDs(ids []int64) ([]User, error)
}

type UserRepo interface {
	CreateUser(*user.User) (int64, error)
	UpdateUser(*user.User) error
	UpdateAvatar(int64, string) error
	DeleteUser(int64) error

	GetByEmail(string) (*user.User, error)
	GetByID(int64) (*user.User, error)
	GetProfileType(int64) (string, error)
	GetAvatarURL(int64) (string, error)

	IsPrivate(int64) (bool, error)
	EmailExists(string, int64) (bool, error)
	UserExists(int64) (bool, error)
	GetByIDs(ids []int64) ([]user.User, error)
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

type UserHandler interface{}
