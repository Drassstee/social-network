package models

import (
	"context"
)

type UserService interface {
	Register(*User) (*UserData, error)
	Login(string, string) (*UserData, error)
	Logout(int64) error
	DeleteUser(int64) error

	GetByID(int64) (*User, error)
	GetByEmail(string) (*User, error)
	GetByIDs([]int64) ([]User, error)

	GetProfile(int64, int64) (*Profile, error)
	GetMe(int64) (*UserData, error)
	UpdateProfile(*User) error
	GetUserID(string) (int64, error)

	Follow(Follow) (string, error)
	Unfollow(Follow) error

	RespondToFollowRequest(Follow) error

	UploadAvatar(context.Context, *Avatar) error
	GetAvatar(int64) (string, error)
}

type UserRepo interface {
	CreateUser(*User) (int64, error)
	UpdateUser(*User) error
	UpdateAvatar(int64, string) error
	DeleteUser(int64) error

	GetByEmail(string) (*User, error)
	GetByID(context.Context, int64) (*User, error)
	GetProfileType(int64) (string, error)
	GetAvatarURL(int64) (string, error)

	IsPrivate(int64) (bool, error)
	EmailExists(string, int64) (bool, error)
	UserExists(int64) (bool, error)
	GetByIDs(ids []int64) ([]User, error)
}

type SessionRepo interface {
	CreateSession(*Session) error
	DeleteSession(string) error
	DeleteAllSessions(int64) error

	GetUserID(string) (int64, error)
	GetUUID(int64) (string, error)
}

type FollowRepo interface {
	CreateFollow(Follow) error
	DeleteFollow(Follow) error
	UpdateFollow(Follow) error

	GetFollowers(int64, string) ([]UserData, error)
	GetFollowing(int64, string) ([]UserData, error)

	IsFollower(int64, int64) (bool, error)
	FollowExists(int64, int64, string) (bool, error)
	GetFollowStatus(int64, int64) (string, error)
}

type UserHandler interface{}
