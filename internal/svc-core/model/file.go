package model

import (
	"mime/multipart"

	"github.com/konsultin/project-goes-here/dto"
)

type File struct {
	Type     dto.FileType
	FileName string
	Size     int64
	Content  multipart.File // Or io.Reader
}
