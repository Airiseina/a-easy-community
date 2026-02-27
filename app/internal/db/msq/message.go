package msq

import (
	"commmunity/app/internal/model"
	"commmunity/app/zlog"

	"go.uber.org/zap"
)

func (db Gorm) SaveMessage(formUserId uint, toUserId uint, content string, tp int) {
	chatMsg := model.Message{
		FromUserID: formUserId,
		ToUserID:   toUserId,
		Content:    content,
		Type:       tp,
	}
	err := db.db.Create(&chatMsg).Error
	if err != nil {
		zlog.Error("保存消息失败", zap.Error(err))
		return
	}
}

func (db Gorm) GetHistoryMessage(userId1 uint, userId2 uint, offset int, limit int) ([]model.Message, error) {
	var chatMsgs []model.Message
	err := db.db.Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", userId1, userId2, userId2, userId1).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&chatMsgs).Error
	if err != nil {
		zlog.Error("查找历史聊天失败", zap.Error(err))
		return nil, err
	}
	return chatMsgs, nil
}

func (db Gorm) SaveNotice(userId uint, senderId uint, typ int, content string, postId uint) {
	notice := model.Notice{
		UserID:   userId,
		SenderID: senderId,
		PostID:   postId,
		Type:     typ,
		Content:  content,
	}
	err := db.db.Create(&notice).Error
	if err != nil {
		zlog.Error("保存通知失败", zap.Error(err))
		return
	}
}

func (db Gorm) GetUnreadNotices(userID uint, offset int, limit int) ([]model.Notice, error) {
	var notices []model.Notice
	err := db.db.Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&notices).Error
	if err != nil {
		zlog.Error("查找未读通知失败", zap.Error(err))
		return nil, err
	}
	return notices, err
}

func (db Gorm) ReadAllNotices(userID uint) error {
	err := db.db.Model(&model.Notice{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
	if err != nil {
		zlog.Error("标记已读失败", zap.Error(err))
		return err
	}
	return nil
}
