package user

import (
	"context"
	"fmt"

	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|


func (us *UserService) GetAvatar(userID int64) (string, error) {
	if userID < 1 {
		return "", fmt.Errorf("get avatar: %w: incorrect user id", models.ErrInvalidData)
	}

	exists, err := us.users.UserExists(userID)
	if err != nil {
		return "", fmt.Errorf("get avatar: %w", err)
	}

	if !exists {
		return "", fmt.Errorf("get avatar: %w: user not found", models.ErrNotFound)
	}

	avatarURL, err := us.users.GetAvatarURL(userID)
	if err != nil {
		return "", fmt.Errorf("get avatar: %w", err)
	}

	return avatarURL, nil
}

//--------------------------------------------------------------------------------------|

func (us *UserService) UploadAvatar(ctx context.Context, a *models.Avatar) error {
	filePath, err := us.uploader.UploadImage(ctx, a.UserID, a.Header.Filename, a.File)
	if err != nil {
		return fmt.Errorf("upload avatar: %w", err)
	}

	err = us.users.UpdateAvatar(a.UserID, filePath)
	if err != nil {
		return fmt.Errorf("upload avatar: %w", err)
	}
	return nil
}
