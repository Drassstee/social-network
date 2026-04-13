package user

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"social-network/internal/models"
	"social-network/internal/models/avatar"
)

// --------------------------------------------------------------------|

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

// --------------------------------------------------------------------|

func (us *UserService) UploadAvatar(a *avatar.Avatar) error {
	ext := filepath.Ext(a.Header.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return fmt.Errorf("upload avatar: %w: invalid file type", models.ErrInvalidData)
	}

	dir := fmt.Sprintf("uploads/avatars/%d", a.UserID)
	os.MkdirAll(dir, 0750)

	filePath := fmt.Sprintf("%s/avatar%s", dir, ext)
	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("upload avatar: %w", err)
	}
	defer dst.Close()
	io.Copy(dst, a.File)

	err = us.users.UpdateAvatar(a.UserID, filePath)
	if err != nil {
		return fmt.Errorf("upload avatar: %w", err)
	}
	return nil
}
