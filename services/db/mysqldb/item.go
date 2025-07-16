package mysqldb

import (
	"log"

	"github.com/darkard2003/wormhole/models"
)

func (r *MySqlRepo) CreateTextItem(item *models.TextItem) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return -1, ToDBError(err, "items", "id")
	}

	defer RecoverDB(tx, &err)

	var id int
	err = tx.QueryRow("INSERT INTO items (user_id, channel_id, type, title, uploaded_at, salt, iv, encryption_metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", item.UserID, item.ChannelID, item.Type, item.Title, item.UploadedAt, item.Salt, item.IV, item.EncryptionMetadata).Scan(&id)

	if err != nil {
		return -1, ToDBError(err, "items", "id")
	}

	_, err = tx.Exec("INSERT INTO text_items (item_id, content) VALUES (?, ?)", id, item.Content)

	if err != nil {
		return -1, ToDBError(err, "items", "id")
	}

	return id, nil
}

func (r *MySqlRepo) CreateFileItem(item *models.FileItem) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return -1, ToDBError(err, "items", "id")
	}

	defer RecoverDB(tx, &err)

	var id int
	err = tx.QueryRow("INSERT INTO items (user_id, channel_id, type, title, uploaded_at, salt, iv, encryption_metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", item.UserID, item.ChannelID, item.Type, item.Title, item.UploadedAt, item.Salt, item.IV, item.EncryptionMetadata).Scan(&id)

	if err != nil {
		return -1, ToDBError(err, "items", "id")
	}

	_, err = tx.Exec("INSERT INTO file_items (item_id, file_name, file_size, mime_type, file_created_at, file_updated_at, file_url) VALUES (?, ?, ?, ?, ?, ?, ?)", id, item.FileName, item.FileSize, item.MimeType, item.FileCreatedAt, item.FileUpdatedAt, item.FileUrl)

	if err != nil {
		return -1, ToDBError(err, "items", "id")
	}

	return id, nil
}
