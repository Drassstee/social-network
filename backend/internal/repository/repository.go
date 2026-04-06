package repository

import (
	"database/sql"

	"social-network/internal/models"
	chatrepo "social-network/internal/repository/chat"
	grouprepo "social-network/internal/repository/group"
	"social-network/internal/repository/notifications"
	"social-network/internal/repository/post"
	"social-network/internal/repository/user"
)

type Repository struct {
	User          models.UserRepo
	Post          models.PostRepo
	Group         models.GroupRepo
	Chat          models.ChatRepo
	Notifications models.NotificationRepo
	DB            *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		User:          user.NewUserRepository(db),
		Post:          post.NewPostRepository(db),
		Group:         grouprepo.NewGroupRepository(db),
		Chat:          chatrepo.NewChatRepository(db),
		Notifications: notifications.NewRepository(db),
		DB:            db,
	}
}
