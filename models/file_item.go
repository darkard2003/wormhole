package models

type FileItem struct {
	Item
	FileName      string
	FileSize      int64
	MimeType      string
	FileCreatedAt string
	FileUpdatedAt string
	FileUrl       string
}
