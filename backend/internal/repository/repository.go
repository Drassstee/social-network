package repository

import (
	"database/sql"

	"social-network/internal/models"
	chatrepo "social-network/internal/repository/chat"
	"social-network/internal/repository/follow"
	grouprepo "social-network/internal/repository/group"
	"social-network/internal/repository/notifications"
	postrepo "social-network/internal/repository/post"
	"social-network/internal/repository/session"
	userrepo "social-network/internal/repository/user"
)

type Repository struct {
	User          models.UserRepo
	Session       *session.SessionRepo
	Follow        *follow.FollowRepo
	Post          models.PostRepo
	Group         models.GroupRepo
	Chat          models.ChatRepo
	Notifications models.NotificationRepo
	DB            *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		User:          userrepo.NewUserRepo(db),
		Session:       session.NewSessionRepo(db),
		Follow:        follow.NewSessionRepo(db),
		Post:          postrepo.NewPostRepo(db),
		Group:         grouprepo.NewGroupRepository(db),
		Chat:          chatrepo.NewChatRepository(db),
		Notifications: notifications.NewRepository(db),
		DB:            db,
	}
}
