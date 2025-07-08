package models

type Item struct {
	ID                 int    `json:"id"`
	UserID             int    `json:"user_id"`
	Type               string `json:"type"`
	Title              string `json:"title"`
	UploadedAt         string `json:"uploaded_at"`
	Salt               string `json:"salt"`
	IV                 string `json:"iv"`
	EncryptionMetadata string `json:"encryption_metadata"`
}
