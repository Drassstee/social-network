package chatsvc

import (
	"context"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

type ChatRepository interface {
	SaveMessage(ctx context.Context, senderID, receiverID int64, body string, imageURL *string) (*models.Message, error)
	GetMessages(ctx context.Context, user1ID, user2ID int64, limit, offset int) ([]models.Message, error)
	MarkAsRead(ctx context.Context, senderID, receiverID int64) error
	GetUnreadCounts(ctx context.Context, userID int64) (map[int64]int, error)
}

//--------------------------------------------------------------------------------------|

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

//--------------------------------------------------------------------------------------|

type FollowRepository interface {
	IsFollower(followerID, followingID int64) (bool, error)
}
