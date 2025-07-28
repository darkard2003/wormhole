package db

import "github.com/darkard2003/wormhole/models"

type DBInterface interface {
	Initialize() error
	Migrate() error

	CreateUser(username, password string, email *string) (int, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	GetAllUsers() ([]*models.User, error)

	CreateChannel(userId int, channelName, channelDescription string, protected bool, password string) (int, error)
	GetChannelById(id int) (*models.Channel, error)
	GetChannelByName(userId int, channelName string) (*models.Channel, error)
	GetChannelsByUserId(userId int) ([]*models.Channel, error)
	UpdateChannel(channel *models.Channel) error
	DeleteChannel(id, userId int) error

	CreateTextItem(item *models.TextItem) (int, error)
	CreateFileItem(item *models.FileItem) (int, error)
	CreateItem(item any) (int, error)
	GetItemById(id int) (any, error)
	PopLatestItem(channelId int) (any, error)
	GetLatestItem(channelId int) (any, error)
}
