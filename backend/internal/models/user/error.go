package user

import "errors"

var (
	ErrExists   = errors.New("email already exists")
	ErrNotFound = errors.New("user not found")

	ErrIncorrectID    = errors.New("incorrect user id")
	ErrIncorrectEmail = errors.New("incorrect email")
	ErrInvalidData    = errors.New("invalid email or password")

	ErrInvalidFormatDate = errors.New("invalid date format, expected YYYY-MM-DD")
	ErrInvalidBirthDate  = errors.New("the date of birth is incorrect")

	ErrPasswordEmpty  = errors.New("password is empty")
	ErrEmailEmpty     = errors.New("email is empty")
	ErrFirstNameEmpty = errors.New("first name is empty")
	ErrLastNameEmpty  = errors.New("last name is empty")
)
