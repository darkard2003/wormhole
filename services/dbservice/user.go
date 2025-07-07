package dbservice

import (
	"database/sql"
	"fmt"

	"github.com/darkard2003/wormhole/models"
)

func (s *DBService) CreateUser(user *models.User) error {
	tx, err := s.DB.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
	}
	_, err = tx.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error inserting user:", err)
		return err
	}
	tx.Commit()
	return nil
}

func (s *DBService) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := s.DB.QueryRow("SELECT id, username, email, password FROM users WHERE username = ?", username).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println("Error querying user:", err)
		return nil, err
	}
	return user, nil
}

func (s *DBService) GetUserById(id int) (*models.User, error) {
	user := &models.User{}
	err := s.DB.QueryRow("SELECT id, username, email, password FROM users WHERE id = ?", id).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println("Error querying user:", err)
		return nil, err
	}
	return user, nil
}

func (s *DBService) UpdateUser(user *models.User) error {
	tx, err := s.DB.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		return err
	}
	_, err = tx.Exec("UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?", user.Username, user.Email, user.Password, user.Id)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error updating user:", err)
		return err
	}
	tx.Commit()
	return nil
}

func (s *DBService) DeleteUser(id int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		fmt.Println("Error starting transaction", err)
		return err
	}

	_, err = tx.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error deleting user:", err)
		return err
	}
	tx.Commit()
	return nil
}

func (s *DBService) GetAllUsers() ([]*models.User, error) {
	rows, err := s.DB.Query("SELECT id, username, email, password FROM users")
	if err != nil {
		fmt.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()
	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		if err != nil {
			fmt.Println("Error scanning user:", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error iterating over users:", err)
		return nil, err
	}
	return users, nil
}
