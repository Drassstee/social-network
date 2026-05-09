package chatsvc

import (
	"context"
	"strings"

	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

type ChatService struct {
	repo       ChatRepository
	userRepo   models.UserRepo
	followRepo models.FollowRepo
}

func NewChatService(repo ChatRepository, userRepo models.UserRepo, followRepo models.FollowRepo) *ChatService {
	return &ChatService{
		repo:       repo,
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

//--------------------------------------------------------------------------------------|

func (s *ChatService) SendMessage(ctx context.Context, senderID, receiverID int64, body string, imageURL *string) (*models.Message, error) {
	return s.repo.SaveMessage(ctx, senderID, receiverID, body, imageURL)
}

//--------------------------------------------------------------------------------------|

func (s *ChatService) GetChatHistory(ctx context.Context, user1ID, user2ID int64, limit, offset int) ([]models.Message, error) {
	return s.repo.GetMessages(ctx, user1ID, user2ID, limit, offset)
}

//--------------------------------------------------------------------------------------|

func (s *ChatService) GetChatableOnlineUsers(ctx context.Context, userID int64, onlineIDs []int64) ([]models.OnlineUser, error) {
	if len(onlineIDs) == 0 {
		return []models.OnlineUser{}, nil
	}

	// 1. Fetch all followers and following for the current user
	followers, err := s.followRepo.GetFollowers(userID, "accept")
	if err != nil {
		return nil, err
	}
	following, err := s.followRepo.GetFollowing(userID, "accept")
	if err != nil {
		return nil, err
	}

	// 2. Map IDs for fast lookup
	canChat := make(map[int64]bool)
	for _, f := range followers {
		canChat[f.ID] = true
	}
	for _, f := range following {
		canChat[f.ID] = true
	}

	// 3. Filter online users
	users, err := s.userRepo.GetByIDs(onlineIDs)
	if err != nil {
		return nil, err
	}

	onlineUsers := make([]models.OnlineUser, 0)
	for _, u := range users {
		if u.ID == userID {
			continue
		}
		if !canChat[u.ID] {
			continue
		}

		name := u.FirstName + " " + u.LastName
		if strings.TrimSpace(name) == "" {
			name = u.Nickname
		}
		onlineUsers = append(onlineUsers, models.OnlineUser{
			ID:       u.ID,
			Username: name,
		})
	}

	return onlineUsers, nil
}
