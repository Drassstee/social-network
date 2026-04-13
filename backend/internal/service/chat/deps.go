package chatsvc

import (
	"context"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

// ChatRepository defines the interface for persisting and retrieving chat messages.
type ChatRepository interface {
	// SaveMessage stores a new chat message in the database.
	SaveMessage(ctx context.Context, senderID, receiverID int, body string, imageURL *string) (*models.Message, error)
	// GetMessages retrieves the chat history between two users, ordered by date.
	GetMessages(ctx context.Context, user1ID, user2ID, limit, offset int) ([]models.Message, error)
}

//--------------------------------------------------------------------------------------|

// UserRepository defines the interface for fetching user data needed by the Hub.
type UserRepository interface {
	GetByID(id int64) (*models.User, error)
}
