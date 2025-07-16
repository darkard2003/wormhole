package mysqldb

import (
	"database/sql"
	"log"

	"github.com/darkard2003/wormhole/models"
	"github.com/darkard2003/wormhole/services/db"
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
	err = tx.QueryRow("INSERT INTO items (user_id, channel_id, type, title, uploaded_at, salt, iv, encryption_metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?) returning id", item.UserID, item.ChannelID, item.Type, item.Title, item.UploadedAt, item.Salt, item.IV, item.EncryptionMetadata).Scan(&id)

	if err != nil {
		return -1, ToDBError(err, "items", "id")
	}

	_, err = tx.Exec("INSERT INTO file_items (item_id, file_name, file_size, mime_type, file_created_at, file_updated_at, file_url) VALUES (?, ?, ?, ?, ?, ?, ?)", id, item.FileName, item.FileSize, item.MimeType, item.FileCreatedAt, item.FileUpdatedAt, item.FileUrl)

	if err != nil {
		return -1, ToDBError(err, "items", "id")
	}

	return id, nil
}

func (r *MySqlRepo) CreateItem(item any) (int, error) {
	switch item := item.(type) {
	case *models.TextItem:
		return r.CreateTextItem(item)
	case *models.FileItem:
		return r.CreateFileItem(item)
	default:
		return -1, db.NewValidationError("item", "type")
	}
}

func (r *MySqlRepo) GetItemById(id int) (any, error) {
	item := &models.Item{}
	textItem := &models.TextItem{}
	fileItem := &models.FileItem{}

	err := r.DB.QueryRow("SELECT id, user_id, channel_id, type, title, uploaded_at, salt, iv, encryption_metadata FROM items WHERE id = ?", id).Scan(&item.ID, &item.UserID, &item.ChannelID, &item.Type, &item.Title, &item.UploadedAt, &item.Salt, &item.IV, &item.EncryptionMetadata)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, ToDBError(err, "items", "id")
	}

	switch item.Type {
	case models.ItemTypeText:
		textItem.Item = *item
		err = r.DB.QueryRow("SELECT content FROM text_items WHERE item_id = ?", id).Scan(&textItem.Content)
		if err != nil {
			return nil, ToDBError(err, "items", "id")
		}
		return textItem, nil
	case models.ItemTypeFile:
		fileItem.Item = *item
		err = r.DB.QueryRow("SELECT file_name, file_size, mime_type, file_created_at, file_updated_at, file_url FROM file_items WHERE item_id = ?", id).Scan(&fileItem.FileName, &fileItem.FileSize, &fileItem.MimeType, &fileItem.FileCreatedAt, &fileItem.FileUpdatedAt, &fileItem.FileUrl)
		if err != nil {
			return nil, ToDBError(err, "items", "id")
		}
		return fileItem, nil
	}

	return nil, db.NewValidationError("item", "type")
}

func (r *MySqlRepo) PopLatestItem(channelId int) (any, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return nil, ToDBError(err, "items", "id")
	}
	defer RecoverDB(tx, &err)
	item := &models.Item{}
	err = tx.QueryRow("SELECT id, user_id, channel_id, type, title, uploaded_at, salt, iv, encryption_metadata FROM items WHERE channel_id = ? ORDER BY uploaded_at DESC LIMIT 1", channelId).Scan(&item.ID, &item.UserID, &item.ChannelID, &item.Type, &item.Title, &item.UploadedAt, &item.Salt, &item.IV, &item.EncryptionMetadata)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, ToDBError(err, "items", "id")
	}

	textItem := &models.TextItem{}
	fileItem := &models.FileItem{}

	var returnItem any

	switch item.Type {
	case models.ItemTypeText:
		textItem.Item = *item
		err = tx.QueryRow("SELECT content FROM text_items WHERE item_id = ?", item.ID).Scan(&textItem.Content)
		if err != nil {
			return nil, ToDBError(err, "items", "id")
		}
		returnItem = textItem
	case models.ItemTypeFile:
		fileItem.Item = *item
		err = tx.QueryRow("SELECT file_name, file_size, mime_type, file_created_at, file_updated_at, file_url FROM file_items WHERE item_id = ?", item.ID).Scan(&fileItem.FileName, &fileItem.FileSize, &fileItem.MimeType, &fileItem.FileCreatedAt, &fileItem.FileUpdatedAt, &fileItem.FileUrl)
		if err != nil {
			return nil, ToDBError(err, "items", "id")
		}
		returnItem = fileItem
	default:
		return nil, db.NewValidationError("item", "type")
	}

	_, err = tx.Exec("DELETE FROM items WHERE id = ?", item.ID)
	if err != nil {
		return nil, ToDBError(err, "items", "id")
	}
	return returnItem, nil

}

func (r *MySqlRepo) GetLatestItem(channelId int) (any, error) {

	tx, err := r.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return nil, ToDBError(err, "items", "id")
	}
	defer RecoverDB(tx, &err)
	item := &models.Item{}
	err = tx.QueryRow("SELECT id, user_id, channel_id, type, title, uploaded_at, salt, iv, encryption_metadata FROM items WHERE channel_id = ? ORDER BY uploaded_at DESC LIMIT 1", channelId).Scan(&item.ID, &item.UserID, &item.ChannelID, &item.Type, &item.Title, &item.UploadedAt, &item.Salt, &item.IV, &item.EncryptionMetadata)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, ToDBError(err, "items", "id")
	}

	textItem := &models.TextItem{}
	fileItem := &models.FileItem{}

	var returnItem any

	switch item.Type {
	case models.ItemTypeText:
		textItem.Item = *item
		err = tx.QueryRow("SELECT content FROM text_items WHERE item_id = ?", item.ID).Scan(&textItem.Content)
		if err != nil {
			return nil, ToDBError(err, "items", "id")
		}
		returnItem = textItem
	case models.ItemTypeFile:
		fileItem.Item = *item
		err = tx.QueryRow("SELECT file_name, file_size, mime_type, file_created_at, file_updated_at, file_url FROM file_items WHERE item_id = ?", item.ID).Scan(&fileItem.FileName, &fileItem.FileSize, &fileItem.MimeType, &fileItem.FileCreatedAt, &fileItem.FileUpdatedAt, &fileItem.FileUrl)
		if err != nil {
			return nil, ToDBError(err, "items", "id")
		}
		returnItem = fileItem
	default:
		return nil, db.NewValidationError("item", "type")
	}

	return returnItem, nil

}
