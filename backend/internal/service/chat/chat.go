package chatsvc

import (
	"context"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

type ChatService struct {
	repo ChatRepository
}

//--------------------------------------------------------------------------------------|

func NewChatService(repo ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

//--------------------------------------------------------------------------------------|

func (s *ChatService) SendMessage(ctx context.Context, senderID, receiverID int, body string, imageURL *string) (*models.Message, error) {
	return s.repo.SaveMessage(ctx, senderID, receiverID, body, imageURL)
}

//--------------------------------------------------------------------------------------|

func (s *ChatService) GetChatHistory(ctx context.Context, user1ID, user2ID, limit, offset int) ([]models.Message, error) {
	return s.repo.GetMessages(ctx, user1ID, user2ID, limit, offset)
}
