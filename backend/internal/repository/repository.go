package repository

import (
	"database/sql"
	"social-network/internal/repository/user"
)

type Repository struct {
	User *user.UserRepository
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		User: user.NewUserRepository(db),
	}
}
