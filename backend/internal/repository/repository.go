package repository

import (
	"database/sql"

	"social-network/internal/repository/follow"
	"social-network/internal/repository/post"
	"social-network/internal/repository/session"
	"social-network/internal/repository/user"
)

type Repository struct {
	*user.UserRepo
	*session.SessionRepo
	*follow.FollowRepo
	*post.PostRepo
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		UserRepo:    user.NewUserRepo(db),
		SessionRepo: session.NewSessionRepo(db),
		FollowRepo:  follow.NewSessionRepo(db),
		PostRepo:    post.NewPostRepo(db),
	}
}
