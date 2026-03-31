package user

import (
	"net/mail"
	"strings"
	"time"
)

type contextKey string

const Key contextKey = "user_id"

// --------------------------------------------------------------------|

type User struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Password    string    `json:"-"`
	DOB         time.Time `json:"-"`
	DOBStr      string    `json:"dob"`
	Avatar      string    `json:"avatar"`
	Nickname    string    `json:"nickname"`
	AboutMe     string    `json:"about_me"`
	ProfileType string    `json:"profile_type"`
}

type UserData struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UUID      string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
}

// --------------------------------------------------------------------|

func (u *User) ValidateData() error {
	if err := u.isEmpty(); err != nil {
		return err
	}

	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return ErrIncorrectEmail
	}

	const layout = time.DateOnly

	t, err := time.Parse(layout, u.DOBStr)
	if err != nil {
		return ErrInvalidFormatDate
	}

	now := time.Now()

	if t.After(now) {
		return ErrInvalidBirthDate
	}

	minDate := now.AddDate(-122, 0, 0)
	if t.Before(minDate) {
		return ErrInvalidBirthDate
	}
	return nil
}

// --------------------------------------------------------------------|

func (u *User) isEmpty() error {
	if len(strings.TrimSpace(u.Email)) == 0 {
		return ErrEmailEmpty
	}
	if len(strings.TrimSpace(u.FirstName)) == 0 {
		return ErrFirstNameEmpty
	}
	if len(strings.TrimSpace(u.LastName)) == 0 {
		return ErrLastNameEmpty
	}
	return nil
}
