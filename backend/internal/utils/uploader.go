package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

//--------------------------------------------------------------------------------------|

// LocalImageUploader implements handles local file storage for images.
type LocalImageUploader struct {
	BaseDir string
	BaseURL string
}

//--------------------------------------------------------------------------------------|

// NewLocalImageUploader creates a new LocalImageUploader.
func NewLocalImageUploader(baseDir, baseURL string) *LocalImageUploader {
	// Ensure the base directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create upload directory %s: %v\n", baseDir, err)
	}
	return &LocalImageUploader{
		BaseDir: baseDir,
		BaseURL: baseURL,
	}
}

//--------------------------------------------------------------------------------------|

// UploadImage saves the content to local storage and returns the relative URL.
func (u *LocalImageUploader) UploadImage(ctx context.Context, userID int, filename string, content io.Reader) (string, error) {
	// 1. Create a unique filename to avoid collisions
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("u%d_%d%s", userID, time.Now().UnixNano(), ext)

	// 2. Define the path
	path := filepath.Join(u.BaseDir, newFilename)

	// 3. Create the file
	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file on disk: %w", err)
	}
	defer file.Close()

	// 4. Save the content
	if _, err := io.Copy(file, content); err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}

	// 5. Return the URL (e.g., /uploads/image.png)
	return filepath.Join(u.BaseURL, newFilename), nil
}
