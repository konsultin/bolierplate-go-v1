package repository

import (
	"time"

	"github.com/go-konsultin/errk"
	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
)

// GetDownloadFileUrl returns the download URL for a given file path
// expiry is optional, default is 1 hour
func (r *Repository) GetDownloadFileUrl(path string, expiry ...time.Duration) (string, error) {
	// Set default expiry to 1 hour
	exp := 1 * time.Hour
	if len(expiry) > 0 && expiry[0] > 0 {
		exp = expiry[0]
	}

	url, err := r.storage.GetPresignedURL(r.ctx, path, exp)
	if err != nil {
		return "", errk.Trace(err)
	}
	return url, nil
}

// UploadFile uploads a file to storage and returns the DTO with URL
// signature is optional
func (r *Repository) UploadFile(file *model.File, signature ...string) (*dto.File, error) {
	// Determine prefix based on file type from DTO routing
	prefix := file.Type.GetPath()

	// Upload
	// We use UploadFile convenience from storage pkg which handles path generation (timestamp_filename)
	path, err := r.storage.UploadFile(r.ctx, prefix, file.FileName, file.Content, file.Size, nil)
	if err != nil {
		return nil, errk.Trace(err)
	}

	// Get Download URL (Presigned)
	url, err := r.GetDownloadFileUrl(path)
	if err != nil {
		return nil, errk.Trace(err)
	}

	sig := ""
	if len(signature) > 0 {
		sig = signature[0]
	}

	return &dto.File{
		FileName:  file.FileName,
		Url:       url,
		Signature: sig,
	}, nil
}
