package user

import (
	"fmt"
	"time"

	"social-network/internal/models"
	"social-network/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//--------------------------------------------------------------------------------------|

func (s *UserService) Register(u *models.User) (*models.UserData, error) {
	if err := u.ValidateData(); err != nil {
		return nil, fmt.Errorf("%w: %v", models.ErrInvalidData, err)
	}

	exists, err := s.users.EmailExists(u.Email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%w: email already exists", models.ErrConflict)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashed)

	id, err := s.users.CreateUser(u)
	if err != nil {
		return nil, err
	}

	sess := &models.Session{
		ID:        uuid.NewString(),
		UserID:    id,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.sessions.CreateSession(sess); err != nil {
		return nil, err
	}

	return &models.UserData{
		ID:        id,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		AvatarURL: utils.FormatAvatarURL(u.AvatarURL),
		UUID:      sess.ID,
		ExpiresAt: &sess.ExpiresAt,
	}, nil
}

//--------------------------------------------------------------------------------------|

func (s *UserService) Login(email, password string) (*models.UserData, error) {
	u, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", models.ErrInvalidData, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("%w: invalid password", models.ErrInvalidData)
	}

	sess := &models.Session{
		ID:        uuid.NewString(),
		UserID:    u.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.sessions.DeleteAllSessions(u.ID); err != nil {
		return nil, err
	}

	if err := s.sessions.CreateSession(sess); err != nil {
		return nil, err
	}

	return &models.UserData{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		AvatarURL: utils.FormatAvatarURL(u.AvatarURL),
		UUID:      sess.ID,
		ExpiresAt: &sess.ExpiresAt,
	}, nil
}

//--------------------------------------------------------------------------------------|

func (s *UserService) Logout(id int64) error {
	if id < 1 {
		return fmt.Errorf("%w: incorrect user id", models.ErrInvalidData)
	}

	uid, err := s.sessions.GetUUID(id)
	if err != nil {
		return err
	}

	return s.sessions.DeleteSession(uid)
}

//--------------------------------------------------------------------------------------|

func (s *UserService) DeleteUser(id int64) error {
	if id < 1 {
		return fmt.Errorf("%w: incorrect user id", models.ErrInvalidData)
	}

	uid, err := s.sessions.GetUUID(id)
	if err == nil {
		_ = s.sessions.DeleteSession(uid)
	}

	return s.users.DeleteUser(id)
}
