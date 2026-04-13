package avatar

import "mime/multipart"

type Avatar struct {
	File   multipart.File
	Header *multipart.FileHeader
	UserID int64
}
