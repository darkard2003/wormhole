package models

type FileItem struct {
	Item
	FileName      string `json:"file_name"`
	FileSize      int64  `json:"file_size"`
	BlobSize      int64  `json:"blob_size"`
	MimeType      string `json:"mime_type"`
	FileCreatedAt string `json:"file_created_at"`
	FileUpdatedAt string `json:"file_updated_at"`
	FileId        string `json:"file_id"`
}
