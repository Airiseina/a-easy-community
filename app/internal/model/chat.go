package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromUserID uint   `gorm:"index" json:"from_user_id"`
	ToUserID   uint   `gorm:"index" json:"to_user_id"`
	Content    string `gorm:"type:longtext" json:"content"`
	Type       int    `gorm:"type:tinyint;comment 类型 1: 文本,2: 图片" json:"type"`
}

type Notice struct {
	gorm.Model
	UserID   uint   `gorm:"index" json:"user_id"`
	Type     int    `gorm:"type:tinyint;comment 类型 1:点赞, 2:评论" json:"type"`
	SenderID uint   `gorm:"index" json:"sender_id"`
	PostID   uint   `gorm:"index" json:"post_id"`
	Content  string `gorm:"type:longtext" json:"content"`
	IsRead   bool   `gorm:"default:false" json:"is_read"`
}

type MessageRequest struct {
	ToUserID uint   `json:"to_user_id"`
	Content  string `json:"content"`
	Type     int    `json:"type"`
}
