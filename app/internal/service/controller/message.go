package controller

import (
	"commmunity/app/internal/db/global"
	"encoding/json"
	"fmt"
)

type MessageDTO struct {
	FromUserID     uint   `json:"from_user_id"`
	FromUserName   string `json:"from_user_name"`
	FromUserAvatar string `json:"from_user_avatar"`
	ToUserID       uint   `json:"to_user_id"`
	ToUserName     string `json:"to_user_name"`
	ToUserAvatar   string `json:"to_user_avatar"`
	Content        string `json:"content"`
	Type           int    `json:"type"`
}

func GetHistoryMessage(forUserId uint, toUserId uint, offset int, pageSize int) ([]MessageDTO, error) {
	mc, err := global.MessageRedis.GetMessageCache(forUserId, toUserId, offset, pageSize)
	if err != nil {
		return nil, err
	}
	if mc == "[]" {
		return []MessageDTO{}, nil
	}
	if mc != "" {
		var cachedPosts []MessageDTO
		if err = json.Unmarshal([]byte(mc), &cachedPosts); err == nil {
			return cachedPosts, nil
		} else {
			return nil, err
		}
	}
	val, err, _ := requestGroup.Do(fmt.Sprintf("HistoryMessage:%d:%d:%d;%d", forUserId, toUserId, offset, pageSize), func() (interface{}, error) {
		ms, err := global.Message.GetHistoryMessage(forUserId, toUserId, offset, pageSize)
		if err != nil {
			return nil, err
		}
		forUser, err := global.Post.GetUserProfile(forUserId)
		if err != nil {
			return nil, err
		}
		forUserName := forUser.UserProfile.Name
		forUserAvatar := forUser.UserProfile.Avatar
		toUser, err := global.Post.GetUserProfile(toUserId)
		if err != nil {
			return nil, err
		}
		toUserName := toUser.UserProfile.Name
		toUserAvatar := toUser.UserProfile.Avatar
		messages := make([]MessageDTO, len(ms))
		for i, m := range ms {
			var fromName, fromAvatar, toName, toAvatar string
			if m.FromUserID == forUserId {
				fromName = forUserName
				fromAvatar = forUserAvatar
				toName = toUserName
				toAvatar = toUserAvatar
			} else {
				fromName = toUserName
				fromAvatar = toUserAvatar
				toName = forUserName
				toAvatar = forUserAvatar
			}
			messages[i] = MessageDTO{
				FromUserID:     m.FromUserID,
				FromUserName:   fromName,
				FromUserAvatar: fromAvatar,
				ToUserID:       m.ToUserID,
				ToUserName:     toName,
				ToUserAvatar:   toAvatar,
				Content:        m.Content,
				Type:           m.Type,
			}
		}
		if len(messages) == 0 {
			err = global.MessageRedis.SetMessageCache(forUserId, toUserId, []MessageDTO{}, offset, pageSize)
			if err != nil {
				return nil, err
			}
			return []MessageDTO{}, nil
		}
		err = global.MessageRedis.SetMessageCache(forUserId, toUserId, messages, offset, pageSize)
		if err != nil {
			return nil, err
		}
		return messages, nil
	})
	if err != nil {
		return nil, err
	}
	return val.([]MessageDTO), nil
}

type NoticeDTO struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	Type      int    `json:"type"`
	SenderID  uint   `json:"sender_id"`
	PostID    uint   `json:"post_id"`
	CreatedAt string `json:"created_at"`
	IsRead    bool   `json:"is_read"`
}

func GetUnreadNotice(userId uint, offset int, limit int) ([]NoticeDTO, error) {
	notices, err := global.Message.GetUnreadNotices(userId, offset, limit)
	if err != nil {
		return nil, err
	}
	err = global.Message.ReadAllNotices(userId)
	if err != nil {
		return nil, err
	}
	noticeDTOs := make([]NoticeDTO, len(notices))
	for i, n := range notices {
		noticeDTOs[i] = NoticeDTO{
			ID:        n.ID,
			Content:   n.Content,
			Type:      n.Type,
			SenderID:  n.SenderID,
			PostID:    n.PostID,
			CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
			IsRead:    n.IsRead,
		}
	}
	return noticeDTOs, nil
}
