package mysqldb

import (
	"database/sql"
	"log"

	"github.com/darkard2003/wormhole/internals/models"
)

func (s *MySqlRepo) CreateUser(username, password string, email *string) (int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return -1, ToDBError(err, "users", "username")
	}

	defer RecoverDB(tx, &err)

	var id int
	err = tx.QueryRow("INSERT INTO users (username, email, password) VALUES (?, ?, ?) RETURNING id", username, email, password).Scan(&id)
	if err != nil {
		return -1, ToDBError(err, "users", "username")
	}

	tx.Commit()
	return id, nil
}

func (s *MySqlRepo) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}

	err := s.DB.QueryRow("SELECT id, username, email, password FROM users WHERE username = ?", username).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error querying user:", err)
		return nil, ToDBError(err, "users", "username")
	}
	return user, nil
}

func (s *MySqlRepo) GetUserById(id int) (*models.User, error) {
	user := &models.User{}
	err := s.DB.QueryRow("SELECT id, username, email, password FROM users WHERE id = ?", id).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error querying user:", err)
		return nil, ToDBError(err, "users", "id")
	}
	return user, nil
}

func (s *MySqlRepo) UpdateUser(user *models.User) error {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return ToDBError(err, "users", "id")
	}
	defer RecoverDB(tx, &err)
	_, err = tx.Exec("UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?", user.Username, user.Email, user.Password, user.Id)
	if err != nil {
		tx.Rollback()
		log.Println("Error updating user:", err)
		return ToDBError(err, "users", "id")
	}
	tx.Commit()
	return nil
}

func (s *MySqlRepo) DeleteUser(id int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction", err)
		return ToDBError(err, "users", "id")
	}

	defer RecoverDB(tx, &err)
	_, err = tx.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Println("Error deleting user:", err)
		return ToDBError(err, "users", "id")
	}
	tx.Commit()
	return nil
}

func (s *MySqlRepo) GetAllUsers() ([]*models.User, error) {
	rows, err := s.DB.Query("SELECT id, username, email, password FROM users")
	if err != nil {
		log.Println("Error querying users:", err)
		return nil, ToDBError(err, "users", "id")
	}
	defer rows.Close()
	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		if err != nil {
			log.Println("Error scanning user:", err)
			return nil, ToDBError(err, "users", "id")
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error iterating over users:", err)
		return nil, ToDBError(err, "users", "id")
	}
	return users, nil
}
