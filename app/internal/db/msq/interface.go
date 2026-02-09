package msq

import "commmunity/app/internal/model"

type UserData interface {
	CreateUser(account string, hash string, name string) error
	GetUser(account string) (*model.User, error)
	GetProfile(account string) (model.User, error)
	DeleteUser(account string) error
	ChangePassword(account string, hash string) error
	ChangeUserName(account string, name string) error
	ChangeAvatar(account string, avatar string) error
	ChangeIntroduction(account string, introduction string) error
	GetUserId(account string) (model.User, error)
	Muted(userID uint, isMuted bool) error
}

type PostData interface {
	CreatePost(userID uint, title string, content string) error
	GetPostList(offset int, pageSize int) ([]model.Post, error)
	GetPostDetail(postID uint) (model.Post, error)
	CreateComment(userID uint, postID uint, content string) error
	GetUserProfile(userID uint) (model.User, error)
	DeletePost(postID uint) error
	DeleteComment(commentID uint) error
	GetCommentDetail(commentID uint) (model.Comment, error)
}
