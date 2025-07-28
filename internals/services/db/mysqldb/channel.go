package mysqldb

import (
	"database/sql"
	"log"

	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/darkard2003/wormhole/models"
)

func (s *MySqlRepo) CreateChannel(userId int, channelName, channelDescription string, protected bool, password string) (int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return -1, ToDBError(err, "channels", "name")
	}

	defer RecoverDB(tx, &err)

	var userExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userId).Scan(&userExists)
	if err != nil {
		log.Println("Error checking user existence:", err)
		return -1, ToDBError(err, "users", "id")
	}
	if !userExists {
		return -1, db.NewNotFoundError("users", "id")
	}

	var channelId int
	err = tx.QueryRow("INSERT INTO channels (user_id, name, description, protected, password) VALUES (?, ?, ?, ?, ?) RETURNING id", userId, channelName, channelDescription, protected, password).Scan(&channelId)

	if err != nil {
		log.Println("Error inserting channel:", ToDBError(err, "channels", "name"))
		return -1, ToDBError(err, "channels", "name")
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return -1, ToDBError(err, "channels", "name")
	}

	return channelId, nil
}

func (s *MySqlRepo) GetChannelById(id int) (*models.Channel, error) {
	channel := &models.Channel{}
	err := s.DB.QueryRow("SELECT id, user_id, name, description, protected, password FROM channels WHERE id = ?", id).Scan(&channel.ID, &channel.UserID, &channel.Name, &channel.Description, &channel.Protected, &channel.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error querying channel:", err)
		return nil, ToDBError(err, "channels", "id")
	}
	return channel, nil
}

func (s *MySqlRepo) GetChannelByName(userId int, channelName string) (*models.Channel, error) {
	channel := &models.Channel{}
	err := s.DB.QueryRow("SELECT id, user_id, name, description, protected, password FROM channels WHERE name = ?", channelName).Scan(&channel.ID, &channel.UserID, &channel.Name, &channel.Description, &channel.Protected, &channel.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error querying channel:", err)
		return nil, ToDBError(err, "channels", "id")
	}
	return channel, nil
}

func (s *MySqlRepo) GetChannelsByUserId(userId int) ([]*models.Channel, error) {
	channels := []*models.Channel{}
	rows, err := s.DB.Query("SELECT id, user_id, name, description, protected, password FROM channels WHERE user_id = ?", userId)
	if err != nil {
		log.Println("Error querying channels:", err)
		return nil, ToDBError(err, "channels", "id")
	}
	defer rows.Close()
	for rows.Next() {
		channel := &models.Channel{}
		if err := rows.Scan(&channel.ID, &channel.UserID, &channel.Name, &channel.Description, &channel.Protected, &channel.Password); err != nil {
			log.Println("Error scanning channel row:", err)
			return nil, err
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func (s *MySqlRepo) UpdateChannel(channel *models.Channel) error {
	tx, err := s.DB.Begin()

	if err != nil {
		log.Println("Error starting transaction for channel update:", err)
		return err
	}

	defer RecoverDB(tx, &err)

	res, err := tx.Exec("UPDATE channels SET name = ?, description = ?, protected = ?, password = ? WHERE id = ?", channel.Name, channel.Description, channel.Protected, channel.Password, channel.ID)

	if err != nil {
		log.Println("Error updating channel:", err)
		tx.Rollback()
		return ToDBError(err, "channels", "id")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return db.NewNotFoundError("channels", "id")
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction for channel update:", err)
		tx.Rollback()
		return ToDBError(err, "channels", "id")
	}

	return nil
}

func (s *MySqlRepo) DeleteChannel(id, userId int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction for channel deletion:", err)
		return err
	}
	defer RecoverDB(tx, &err)
	res, err := tx.Exec("DELETE FROM channels WHERE id = ? and user_id = ?", id, userId)
	if err != nil {
		log.Println("Error deleting channel:", err)
		tx.Rollback()
		return ToDBError(err, "channels", "id")
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		log.Println("No channel found with the given ID for the user")
		tx.Rollback()
		return db.NewNotFoundError("channels", "id")
	}
	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction for channel deletion:", err)
		tx.Rollback()
		return ToDBError(err, "channels", "id")
	}
	return nil
}
