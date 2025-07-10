package mysqldb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/darkard2003/wormhole/models"
)

func (s *MySqlRepo) CreateChannel(channel *models.Channel, userId int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Query("SELECT COUNT(*) FROM channels WHERE user_id = ? AND name = ?", userId, channel.Name)
	if err != nil {
		log.Println("Error checking channel existence:", err)
		tx.Rollback()
		return err
	}
	var count int
	if res.Next() {
		if err := res.Scan(&count); err != nil {
			log.Println("Error scanning channel count:", err)
			tx.Rollback()
			return err
		}
	}
	res.Close()
	if count > 0 {
		log.Println("Channel with this name already exists for the user")
		tx.Rollback()
		return fmt.Errorf("channel exists")
	}
	_, err = tx.Exec("INSERT INTO channels (user_id, name, description, protected, password) VALUES (?, ?, ?, ?, ?)", userId, channel.Name, channel.Description, channel.Protected, channel.Password)
	if err != nil {
		log.Println("Error inserting channel:", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
	}

	return nil
}

func (s *MySqlRepo) GetChannelById(id int) (*models.Channel, error) {
	channel := &models.Channel{}
	err := s.DB.QueryRow("SELECT id, user_id, name, description, protected, password FROM channels WHERE id = ?", id).Scan(&channel.ID, &channel.UserID, &channel.Name, &channel.Description, &channel.Protected, &channel.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error querying channel:", err)
		return nil, err
	}
	return channel, nil
}

func (s *MySqlRepo) GetChannelsByUserId(userId int) ([]*models.Channel, error) {
	channels := []*models.Channel{}
	rows, err := s.DB.Query("SELECT id, user_id, name, description, protected, password FROM channels WHERE user_id = ?", userId)
	if err != nil {
		log.Println("Error querying channels:", err)
		return nil, err
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
	_, err = tx.Exec("UPDATE channels SET name = ?, description = ?, protected = ?, password = ? WHERE id = ?", channel.Name, channel.Description, channel.Protected, channel.Password, channel.ID)

	if err != nil {
		log.Println("Error updating channel:", err)
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction for channel update:", err)
		tx.Rollback()
		return err
	}

	return nil
}

func (s *MySqlRepo) DeleteChannel(id, userId int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction for channel deletion:", err)
		return err
	}
	res, err := tx.Exec("DELETE FROM channels WHERE id = ? and user_id = ?", id, userId)
	if err != nil {
		log.Println("Error deleting channel:", err)
		tx.Rollback()
		return err
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		log.Println("No channel found with the given ID for the user")
		tx.Rollback()
		return fmt.Errorf("channel not found")
	}
	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction for channel deletion:", err)
		tx.Rollback()
		return err
	}
	return nil
}
