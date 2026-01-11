package dto

type FileType string

const (
	FileType_USER_AVATAR    FileType = "user-avatar"
	FileType_JOB_ATTACHMENT FileType = "job-attachment"
	FileType_DOCUMENT       FileType = "document"
)

type UploadFile_Payload struct {
	Type     FileType `json:"type" validate:"required"`
	FileName string   `json:"fileName" validate:"required"`
	Size     int64    `json:"size" validate:"required"`
	// Content is handled separately (e.g. multipart/form-data) or stream
}

func (t FileType) GetPath() string {
	switch t {
	case FileType_USER_AVATAR:
		return "users/avatar"
	case FileType_JOB_ATTACHMENT:
		return "jobs/attachment"
	case FileType_DOCUMENT:
		return "documents"
	default:
		return string(t)
	}
}
