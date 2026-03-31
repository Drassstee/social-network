package user

import (
	"social-network/internal/models/profile"
	"social-network/internal/models/user"
)

type UserHandler struct {
	Users UserServ
}

type UserServ interface {
	Register(*user.User) (*user.UserData, error)
	Login(string, string) (*user.UserData, error)
	Logout(int64) error

	GetProfile(int64) (*profile.Profile, error)
	UpdateProfile(*user.User) error
	GetUserID(string) (int64, error)
}

func NewUserHandler(serv UserServ) *UserHandler {
	return &UserHandler{Users: serv}
}
