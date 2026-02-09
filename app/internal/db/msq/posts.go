package msq

import (
	"commmunity/app/internal/model"
	"commmunity/app/zlog"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (db Gorm) CreatePost(userID uint, title string, content string) error {
	post := model.Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}
	result := db.db.Create(&post)
	if result.Error != nil {
		zlog.Error("帖子创建失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func (db Gorm) GetPostList(offset int, pageSize int) ([]model.Post, error) {
	var posts []model.Post
	err := db.db.Preload("User"). //预加载用户信息后，文章信息只要id和标题和跟热度有关的，一个网页中仅给规定文章
					Preload("User.UserProfile").
					Select("id, user_id, title, created_at, view_count, like_count, comment_count").
					Order("created_at desc").
					Offset(offset).
					Limit(pageSize).
					Find(&posts).Error
	if err != nil {
		zlog.Error("论坛生成失败", zap.Error(err))
		return nil, err
	}
	for _, post := range posts {
		fmt.Println(post)
	}
	return posts, nil
}

//// 获得作者所有的文章
//func (db Gorm) GetPosts(userID uint) ([]model.Post, error) {}

func (db Gorm) GetPostDetail(postID uint) (model.Post, error) {
	var post model.Post
	err := db.db.Preload("User").
		Preload("User.UserProfile").
		Preload("Comments").
		Preload("Comments.User").
		Preload("Comments.User.UserProfile").
		First(&post, postID).Error
	if err != nil {
		zlog.Error("查询文章失败", zap.Error(err))
		return model.Post{}, err
	}
	return post, nil
}

func (db Gorm) CreateComment(userID uint, postID uint, content string) error {
	comment := model.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}
	result := db.db.Create(&comment)
	if result.Error != nil {
		zlog.Error("评论创建失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func (db Gorm) GetUserProfile(userID uint) (model.User, error) {
	var user model.User
	err := db.db.Preload("UserProfile").Preload("Posts").First(&user, userID).Error
	if err != nil {
		zlog.Error("查找失败", zap.Error(err))
		return model.User{}, err
	}
	return user, nil
}

func (db Gorm) DeletePost(postID uint) error {
	post := model.Post{
		Model: gorm.Model{
			ID: postID,
		},
	}
	return db.db.Select("Comments").Delete(&post).Error
}

func (db Gorm) DeleteComment(commentID uint) error {
	comment := model.Comment{
		Model: gorm.Model{
			ID: commentID,
		},
	}
	return db.db.Delete(&comment).Error
}

func (db Gorm) GetCommentDetail(commentID uint) (model.Comment, error) {
	var comment model.Comment
	err := db.db.Preload("User").Preload("User.UserProfile").First(&comment, commentID).Error
	if err != nil {
		zlog.Error("查询失败", zap.Error(err))
		return model.Comment{}, err
	}
	return comment, nil
}
