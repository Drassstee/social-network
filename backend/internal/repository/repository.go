package repository

import (
	"database/sql"

	"social-network/internal/models"
	repost "social-network/internal/repository/post"
	"social-network/internal/repository/user"
)

type Repository struct {
	User models.UserRepo
	Post models.PostRepo
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		User: user.NewUserRepository(db),
		Post: repost.NewPostRepository(db),
	}
}
