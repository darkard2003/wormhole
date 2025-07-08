package dbservice

import (
	"database/sql"
	"log"

	"github.com/darkard2003/wormhole/models"
)

func (s *DBService) CreateChannel(channel *models.Channel, userId int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
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

func (s *DBService) GetChannelById(id int) (*models.Channel, error) {
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

func (s *DBService) GetChannelsByUserId(userId int) ([]*models.Channel, error) {
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

func (s *DBService) UpdateChannel(channel *models.Channel) error {
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

func (s *DBService) DeleteChannel(id int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction for channel deletion:", err)
		return err
	}
	_, err = tx.Exec("DELETE FROM channels WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting channel:", err)
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction for channel deletion:", err)
		tx.Rollback()
		return err
	}
	return nil
}
