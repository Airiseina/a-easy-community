package msq

import (
	"commmunity/app/internal/model"
	"commmunity/app/zlog"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (db Gorm) CreateUser(account string, hash string, name string) error {
	user := model.User{
		Account: account,
		Hash:    hash,
		UserProfile: model.UserProfile{
			Name: name,
		},
	}
	return db.db.Create(&user).Error
}

func (db Gorm) GetProfile(account string) (model.User, error) {
	var user model.User
	result := db.db.Where("account = ?", account).Preload("UserProfile").Preload("Posts").First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			zlog.Warn("用户未找到")
			return model.User{}, nil
		} else {
			zlog.Error("找用户时出错了", zap.Error(result.Error))
			return model.User{}, result.Error
		}
	}
	return user, nil
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

// 不删评论
func (db Gorm) DeleteUser(account string) error {
	var user model.User
	if err := db.db.Where("account = ?", account).First(&user).Error; err != nil {
		return err
	}
	return db.db.Unscoped().Select("UserProfile", "Posts").Delete(&user).Error
}
func (db Gorm) ChangePassword(account string, hash string) error {
	user := model.User{
		Account: account,
		Hash:    hash,
	}
	result := db.db.Model(&model.User{}).Where("account = ?", account).Select("hash").Updates(user)
	if result.Error != nil {
		zlog.Error("密码修改失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func (db Gorm) ChangeUserName(account string, name string) error {
	user, err := db.GetUserId(account)
	if err != nil {
		return err
	}
	err = db.db.Model(&model.UserProfile{}).Where("user_id = ?", user.ID).Update("name", name).Error
	if err != nil {
		zlog.Error("用户名修改失败", zap.Error(err))
		return err
	}
	return nil
}

func (db Gorm) ChangeAvatar(account string, avatar string) error {
	user, err := db.GetUserId(account)
	if err != nil {
		return err
	}
	err = db.db.Model(&model.UserProfile{}).Where("user_id = ?", user.ID).Update("avatar", avatar).Error
	if err != nil {
		zlog.Error("头像修改失败", zap.Error(err))
		return err
	}
	return nil
}

func (db Gorm) ChangeIntroduction(account string, introduction string) error {
	user, err := db.GetUserId(account)
	if err != nil {
		return err
	}
	err = db.db.Model(&model.UserProfile{}).Where("user_id = ?", user.ID).Update("introduction", introduction).Error
	if err != nil {
		zlog.Error("简介修改失败", zap.Error(err))
		return err
	}
	return nil
}

func (db Gorm) GetUserId(account string) (model.User, error) {
	var user model.User
	err := db.db.Model(&model.User{}).Preload("UserProfile").Select("id").Where("account = ?", account).First(&user).Error
	if err != nil {
		zlog.Error("查找id失败", zap.Error(err))
		return model.User{}, err
	}
	return user, nil
}

func (db Gorm) Muted(userID uint, isMuted bool) error {
	err := db.db.Model(&model.UserProfile{}).Where("user_id = ?", userID).Update("is_muted", isMuted).Error
	if err != nil {
		zlog.Error("禁言失败", zap.Error(err))
		return err
	}
	return nil
}
