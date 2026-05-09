package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

//--------------------------------------------------------------------------------------|

type LocalImageUploader struct {
	BaseDir string
	BaseURL string
}

//--------------------------------------------------------------------------------------|

func NewLocalImageUploader(baseDir, baseURL string) *LocalImageUploader {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create upload directory %s: %v\n", baseDir, err)
	}
	return &LocalImageUploader{
		BaseDir: baseDir,
		BaseURL: baseURL,
	}
}

//--------------------------------------------------------------------------------------|

func (u *LocalImageUploader) UploadImage(ctx context.Context, userID int64, filename string, content io.Reader) (string, error) {
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("u%d_%d%s", userID, time.Now().UnixNano(), ext)

	diskPath := filepath.Join(u.BaseDir, newFilename)

	file, err := os.Create(diskPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file on disk: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, content); err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}
	return path.Join(u.BaseURL, newFilename), nil
}
