package models

import "errors"

//--------------------------------------------------------------------------------------|

// ValidationError represents a validation failure on a specific field.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

//--------------------------------------------------------------------------------------|

// NotFoundError indicates that a requested entity could not be found.
type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return e.Entity + " not found"
}

//--------------------------------------------------------------------------------------|

// AuthorizationError indicates that the user lacks permission for the action.
type AuthorizationError struct {
	UserID int
}

func (e *AuthorizationError) Error() string {
	if e.UserID == 0 {
		return "authentication required"
	}
	return "access denied"
}

//--------------------------------------------------------------------------------------|

// Sentinel errors used across the application.
var (
	// Auth & user errors.
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionNotFound    = errors.New("session not found")

	// Validation errors.
	ErrUsernameLength     = errors.New("username must be between 3 and 20 characters")
	ErrUsernameFormat     = errors.New("username contains invalid characters")
	ErrPasswordLength     = errors.New("password must be at least 8 characters")
	ErrPasswordComplexity = errors.New("password must contain uppercase, lowercase, and digits")
	ErrInvalidEmail       = errors.New("invalid email address")

	// Post errors.
	ErrPostNotFound     = errors.New("post not found")
	ErrInvalidPostTitle = errors.New("post title is required")
	ErrInvalidPostBody  = errors.New("post body is required")

	// Comment errors.
	ErrCommentNotFound       = errors.New("comment not found")
	ErrParentCommentNotFound = errors.New("parent comment not found")
	ErrInvalidCommentBody    = errors.New("comment body is required")
	ErrCommentDepthExceeded  = errors.New("maximum comment depth exceeded")

	// Category errors.
	ErrInvalidCategoryName = errors.New("category name is required")

	// Image errors.
	ErrImageTooBig        = errors.New("image exceeds maximum allowed size")
	ErrInvalidImageFormat = errors.New("unsupported image format")

	// Database constraint errors.
	ErrUniqueConstraint     = errors.New("unique constraint violation")
	ErrForeignKeyConstraint = errors.New("foreign key constraint violation")

	// Group errors.
	ErrGroupNotFound       = errors.New("group not found")
	ErrAlreadyMember       = errors.New("user is already a member of this group")
	ErrNotMember           = errors.New("user is not a member of this group")
	ErrNotGroupCreator     = errors.New("only the group creator can perform this action")
	ErrInvitationNotFound  = errors.New("invitation not found")
	ErrJoinRequestNotFound = errors.New("join request not found")
	ErrEventNotFound       = errors.New("event not found")
	ErrAlreadyInvited      = errors.New("user has already been invited")
	ErrAlreadyRequested    = errors.New("join request already pending")
)
