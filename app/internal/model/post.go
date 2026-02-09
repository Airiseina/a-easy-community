package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title        string    `gorm:"type:varchar(100);not null" json:"title"`
	Content      string    `gorm:"type:longtext" json:"content"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID;not null" json:"user"`
	Comments     []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	ViewCount    uint      `gorm:"default:0" json:"view_count"`    // 浏览量
	LikeCount    uint      `gorm:"default:0" json:"like_count"`    // 点赞数
	CommentCount uint      `gorm:"default:0" json:"comment_count"` // 评论数
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:longtext;not null" json:"content"`
	PostID  uint   `gorm:"index;not null" json:"post_id"`
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	User    User   `gorm:"foreignKey:UserID;not null" json:"user"`
}

type PostRequest struct {
	Title   string
	Content string
}

type CommentRequest struct {
	Content string
}
