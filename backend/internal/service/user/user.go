package user

import (
	"context"
	"social-network/internal/models"
)

type UserService struct {
	repo models.UserRepo
}

func NewUserService(repo models.UserRepo) *UserService {
	return &UserService{repo: repo}
}
func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) GetByIDs(ctx context.Context, ids []int) ([]models.User, error) {
	return s.repo.GetByIDs(ctx, ids)
}
