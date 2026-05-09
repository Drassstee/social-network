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

//--------------------------------------------------------------------------------------|

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

// MarkAsRead marks all messages from sender to receiver as read.
func (s *ChatService) MarkAsRead(ctx context.Context, senderID, receiverID int64) error {
	return s.repo.MarkAsRead(ctx, senderID, receiverID)
}

//--------------------------------------------------------------------------------------|

// GetUnreadCounts returns unread counts for a user.
func (s *ChatService) GetUnreadCounts(ctx context.Context, userID int64) (map[int64]int, error) {
	return s.repo.GetUnreadCounts(ctx, userID)
}

//--------------------------------------------------------------------------------------|

// GetChatableUsers returns all users that the current user can chat with.
func (s *ChatService) GetChatableUsers(ctx context.Context, userID int64) ([]models.OnlineUser, error) {
	// 1. Fetch all followers and following for the current user
	followers, err := s.followRepo.GetFollowers(userID, "accept")
	if err != nil {
		return nil, err
	}
	following, err := s.followRepo.GetFollowing(userID, "accept")
	if err != nil {
		return nil, err
	}

	// 2. Map IDs for fast lookup and deduplication
	ids := make(map[int64]bool)
	for _, f := range followers {
		ids[f.ID] = true
	}
	for _, f := range following {
		ids[f.ID] = true
	}

	if len(ids) == 0 {
		return []models.OnlineUser{}, nil
	}

	idList := make([]int64, 0, len(ids))
	for id := range ids {
		idList = append(idList, id)
	}

	// 3. Fetch user details
	users, err := s.userRepo.GetByIDs(idList)
	if err != nil {
		return nil, err
	}

	chatable := make([]models.OnlineUser, 0, len(users))
	for _, u := range users {
		name := u.FirstName + " " + u.LastName
		if strings.TrimSpace(name) == "" {
			name = u.Nickname
		}
		chatable = append(chatable, models.OnlineUser{
			ID:       u.ID,
			Username: name,
		})
	}

	return chatable, nil
}

