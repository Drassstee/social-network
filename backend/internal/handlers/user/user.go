package user

import (
	"social-network/internal/models/follow"
	"social-network/internal/models/profile"
	"social-network/internal/models/user"
)

var errInternalServer = map[string]string{"error": "internal server error"}

type UserHandler struct {
	Users UserServ
}

type UserServ interface {
	Register(*user.User) (*user.UserData, error)
	Login(string, string) (*user.UserData, error)
	Logout(int64) error

	GetProfile(int64, int64) (*profile.Profile, error)
	UpdateProfile(*user.User) error
	GetUserID(string) (int64, error)

	Follow(follow.Follow) (string, error)
	Unfollow(follow.Follow) error

	GetNotification(int64) ([]user.UserData, error)
	RespondToFollowRequest(follow.Follow) error
}

func NewUserHandler(serv UserServ) *UserHandler {
	return &UserHandler{Users: serv}
}
