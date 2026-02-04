package msq

import (
	"commmunity/app/internal/model"
	"commmunity/app/zlog"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (db Gorm) CreateUser(user model.User) error {
	result := db.db.Create(&user)
	return result.Error
}

func (db Gorm) GetUser(account string) (*model.User, error) {
	var user model.User
	result := db.db.Where("account = ?", account).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			zlog.Warn("用户未找到")
			return nil, nil
		} else {
			zlog.Error("找用户时出错了")
			return nil, result.Error
		}
	}
	return &user, nil
}

func (db Gorm) DeleteUser(account string) error {
	result := db.db.Where("account = ?", account).Delete(&model.User{})
	return result.Error
}

func (db Gorm) ChangePassword(user model.User) error {
	result := db.db.Model(&model.User{}).Where("account = ?", user.Account).Select("hash").Updates(&user)
	if result.Error != nil {
		zlog.Error("密码修改失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func (db Gorm) ChangeUserName(user model.User) error {
	result := db.db.Model(&model.User{}).Where("account = ?", user.Account).Select("Name").Updates(&user)
	if result.Error != nil {
		zlog.Error("用户名修改失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func (db Gorm) ChangeAvatar(user model.User) error {
	result := db.db.Model(&model.User{}).Where("account = ?", user.Account).Select("avatar").Updates(&user)
	if result.Error != nil {
		zlog.Error("头像修改失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func (db Gorm) ChangeIntroduction(user model.User) error {
	result := db.db.Model(&model.User{}).Where("account = ?", user.Account).Select("Introduction").Updates(&user)
	if result.Error != nil {
		zlog.Error("简介修改失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}
