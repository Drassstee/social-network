package user

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"social-network/internal/models"
	"social-network/internal/models/session"
	"social-network/internal/models/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// --------------------------------------------------------------------|

func (us *UserService) Register(u *user.User) (*user.UserData, error) {
	if len(strings.TrimSpace(u.Password)) == 0 {
		return nil, fmt.Errorf("register: %w: password is empty", models.ErrInvalidData)
	}

	err := u.ValidateData()
	if err != nil {
		return nil, fmt.Errorf("register: %w: %w", models.ErrInvalidData, err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(password)

	id, err := us.users.CreateUser(u)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	u.ID = id

	var s = session.Session{
		UserID:    id,
		UUID:      uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = us.sessions.CreateSession(&s)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return &user.UserData{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		UUID:      s.UUID,
		ExpiresAt: &s.ExpiresAt,
	}, nil
}

// --------------------------------------------------------------------|

func (us *UserService) Login(email, password string) (*user.UserData, error) {
	if len(strings.TrimSpace(email)) == 0 {
		return nil, fmt.Errorf("login: %w: email is empty", models.ErrInvalidData)
	}
	if len(strings.TrimSpace(password)) == 0 {
		return nil, fmt.Errorf("login: %w: password is empty", models.ErrInvalidData)
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("login: %w: incorrect email", models.ErrInvalidData)
	}

	u, err := us.users.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return nil, fmt.Errorf("login: %w: invalid email or password", models.ErrInvalidData)
	}

	var s = session.Session{
		UserID:    u.ID,
		UUID:      uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = us.sessions.DeleteAllSessions(u.ID)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	err = us.sessions.CreateSession(&s)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	return &user.UserData{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		UUID:      s.UUID,
		ExpiresAt: &s.ExpiresAt,
	}, nil
}

// --------------------------------------------------------------------|

func (us *UserService) Logout(id int64) error {
	if id < 1 {
		return fmt.Errorf("logout: %w: incorrect user id", models.ErrInvalidData)
	}

	uuid, err := us.sessions.GetUUID(id)
	if err != nil {
		return fmt.Errorf("logout: %w", err)
	}

	err = us.sessions.DeleteSession(uuid)
	if err != nil {
		return fmt.Errorf("logout: %w", err)
	}
	return nil
}

// --------------------------------------------------------------------|
