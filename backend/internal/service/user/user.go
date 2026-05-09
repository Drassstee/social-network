package user

import (
	"context"

	"social-network/internal/models"
)

type UserService struct {
	users         models.UserRepo
	sessions      models.SessionRepo
	follows       models.FollowRepo
	posts         models.PostRepo
	notifications models.NotificationService
	uploader      models.ImageUploader
}

func NewUserService(ur models.UserRepo, sr models.SessionRepo, fr models.FollowRepo, pr models.PostRepo, ns models.NotificationService, uploader models.ImageUploader) *UserService {
	return &UserService{
		users:         ur,
		sessions:      sr,
		follows:       fr,
		posts:         pr,
		notifications: ns,
		uploader:      uploader,
	}
}

func (s *UserService) GetByID(id int64) (*models.User, error) {
	return s.users.GetByID(context.Background(), id)
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	return s.users.GetByEmail(email)
}

func (s *UserService) GetByIDs(ids []int64) ([]models.User, error) {
	return s.users.GetByIDs(ids)
}
