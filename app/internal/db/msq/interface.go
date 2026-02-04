package msq

import "commmunity/app/internal/model"

type UserData interface {
	CreateUser(user model.User) error
	GetUser(account string) (*model.User, error)
	DeleteUser(account string) error
	ChangePassword(user model.User) error
	ChangeUserName(user model.User) error
	ChangeAvatar(user model.User) error
	ChangeIntroduction(user model.User) error
}
