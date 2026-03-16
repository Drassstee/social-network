package repository

import (
	"database/sql"
	"social-network/internal/models"
	"social-network/internal/repository/user"
)

type Repository struct {
	User models.UserRepo
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		User: user.NewUserRepository(db),
	}
}
