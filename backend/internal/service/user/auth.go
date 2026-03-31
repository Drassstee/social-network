package user

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"social-network/internal/models/session"
	"social-network/internal/models/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// --------------------------------------------------------------------|

func (r *UserService) Register(u *user.User) (*user.UserData, error) {
	if len(strings.TrimSpace(u.Password)) == 0 {
		return nil, fmt.Errorf("register: %w", user.ErrPasswordEmpty)
	}

	err := u.ValidateData()
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(password)

	id, err := r.users.CreateUser(u)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err) // err: email already exists || db.err
	}
	u.ID = id

	var s = session.Session{
		UserID:    id,
		UUID:      uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = r.sessions.CreateSession(&s)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err) // err: db.err
	}

	return &user.UserData{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		ExpiresAt: s.ExpiresAt,
	}, nil
}

// --------------------------------------------------------------------|

func (r *UserService) Login(email, password string) (*user.UserData, error) {
	if len(strings.TrimSpace(email)) == 0 {
		return nil, user.ErrEmailEmpty
	}
	if len(strings.TrimSpace(password)) == 0 {
		return nil, user.ErrPasswordEmpty
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, user.ErrIncorrectEmail
	}

	u, err := r.users.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err) // err: db.err || invalid data
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return nil, user.ErrInvalidData // err : invalid data
	}

	var s = session.Session{
		UserID:    u.ID,
		UUID:      uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = r.sessions.CreateSession(&s)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err) // err: db.err
	}

	return &user.UserData{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		ExpiresAt: s.ExpiresAt,
	}, nil
}

// --------------------------------------------------------------------|

func (r *UserService) Logout(id int64) error {
	if id < 1 {
		return user.ErrIncorrectID
	}

	err := r.sessions.Delete(id)
	if err != nil {
		return fmt.Errorf("logout: %w", err) // err: db.err
	}
	return nil
}

// --------------------------------------------------------------------|
