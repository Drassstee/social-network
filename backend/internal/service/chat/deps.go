package chatsvc

import (
	"context"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

// ChatRepository defines the interface for persisting and retrieving chat messages.
type ChatRepository interface {
	// SaveMessage stores a new chat message in the database.
	SaveMessage(ctx context.Context, senderID, receiverID int64, body string, imageURL *string) (*models.Message, error)
	// GetMessages retrieves the chat history between two users, ordered by date.
	GetMessages(ctx context.Context, user1ID, user2ID int64, limit, offset int) ([]models.Message, error)
	// MarkAsRead marks all messages from sender to receiver as read.
	MarkAsRead(ctx context.Context, senderID, receiverID int64) error
	// GetUnreadCounts returns a map of senderID -> unread count for the given user.
	GetUnreadCounts(ctx context.Context, userID int64) (map[int64]int, error)
}

//--------------------------------------------------------------------------------------|

// UserRepository defines the interface for fetching user data needed by the Hub.
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

//--------------------------------------------------------------------------------------|

// FollowRepository defines the interface for checking follow status.
type FollowRepository interface {
	IsFollower(followerID, followingID int64) (bool, error)
}
