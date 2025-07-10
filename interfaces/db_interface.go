package interfaces

import "github.com/darkard2003/wormhole/models"

type DBInterface interface {
	Initialize() error
	Migrate() error
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	GetAllUsers() ([]*models.User, error)
	CreateChannel(channel *models.Channel, userId int) error
	GetChannelById(id int) (*models.Channel, error)
	GetChannelsByUserId(userId int) ([]*models.Channel, error)
	UpdateChannel(channel *models.Channel) error
	DeleteChannel(id, userId int) error
}
