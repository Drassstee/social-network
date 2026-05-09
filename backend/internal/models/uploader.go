package models

import (
	"context"
	"io"
)

//--------------------------------------------------------------------------------------|

type ImageUploader interface {
	UploadImage(ctx context.Context, userID int64, filename string, content io.Reader) (string, error)
}
