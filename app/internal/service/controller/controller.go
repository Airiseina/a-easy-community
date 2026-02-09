package controller

import (
	"commmunity/app/internal/db/global"
	"commmunity/app/internal/model"
)

func CreatePost(account string, title string, content string) (error, bool) {
	user, err := global.User.GetUserId(account)
	if err != nil {
		return err, false
	}
	if user.UserProfile.IsMuted {
		return nil, false
	}
	return global.Post.CreatePost(user.ID, title, content), true
}

type PostsDTO struct {
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	PostID       uint   `json:"post_id"`
	Title        string `json:"title"`
	ViewCount    uint   `json:"view_count"`
	LikeCount    uint   `json:"like_count"`
	CommentCount uint   `json:"comment_count"`
}

func GetPostList(offset int, pageSize int) ([]PostsDTO, error) {
	ps, err := global.Post.GetPostList(offset, pageSize)
	if err != nil {
		return nil, err
	}
	posts := make([]PostsDTO, len(ps))
	for i, p := range ps {
		posts[i] = PostsDTO{
			Name:         p.User.UserProfile.Name,
			Avatar:       p.User.UserProfile.Avatar,
			PostID:       p.ID,
			Title:        p.Title,
			ViewCount:    p.ViewCount,
			LikeCount:    p.LikeCount,
			CommentCount: p.CommentCount,
		}
	}
	return posts, nil
}

type PostDTO struct {
	PostsDTO
	Content string       `json:"content"`
	UserId  uint         `json:"user_id"`
	Comment []CommentDTO `json:"comment"`
}

type CommentDTO struct {
	ID         uint   `json:"id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
	UserId     uint   `json:"user_id"`
}

func GetPostDetail(postId uint) (PostDTO, error) {
	p, err := global.Post.GetPostDetail(postId)
	if err != nil {
		return PostDTO{}, err
	}
	commentDTOs := make([]CommentDTO, 0, len(p.Comments))
	for _, c := range p.Comments {
		commentDTOs = append(commentDTOs, CommentDTO{
			ID:         c.ID,
			Content:    c.Content,
			CreatedAt:  c.CreatedAt.Format("2006-01-02 15:04:05"),
			UserName:   c.User.UserProfile.Name,
			UserAvatar: c.User.UserProfile.Avatar,
			UserId:     c.UserID,
		})
	}
	post := PostDTO{
		PostsDTO: PostsDTO{
			Name:         p.User.UserProfile.Name,
			Avatar:       p.User.UserProfile.Avatar,
			PostID:       p.ID,
			Title:        p.Title,
			ViewCount:    p.ViewCount,
			LikeCount:    p.LikeCount,
			CommentCount: p.CommentCount,
		},
		Content: p.Content,
		UserId:  p.User.ID,
		Comment: commentDTOs,
	}
	return post, nil
}

func CreateComment(account string, postID uint, content string) (error, bool) {
	user, err := global.User.GetUserId(account)
	if err != nil {
		return err, false
	}
	if user.UserProfile.IsMuted {
		return nil, false
	}
	return global.Post.CreateComment(user.ID, postID, content), true
}

type UserProfileDTO struct {
	Account      string        `json:"account"`
	Name         string        `json:"name"`
	Introduction string        `json:"introduction"`
	Avatar       string        `json:"avatar"`
	Role         int           `json:"role"`
	IsMuted      bool          `json:"isMuted"`
	Posts        []UserPostDTO `json:"posts"`
}
type UserPostDTO struct {
	PostID       uint   `json:"post_id"`
	Title        string `json:"title"`
	ViewCount    uint   `json:"view_count"`
	LikeCount    uint   `json:"like_count"`
	CommentCount uint   `json:"comment_count"`
}

func GetUserProfile(id uint) (*UserProfileDTO, error) {
	user, err := global.Post.GetUserProfile(id)
	if err != nil {
		return nil, err
	}
	userPostDTO := make([]UserPostDTO, 0, len(user.Posts))
	for _, post := range user.Posts {
		userPostDTO = append(userPostDTO, UserPostDTO{
			PostID:       post.ID,
			Title:        post.Title,
			ViewCount:    post.ViewCount,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
		})
	}
	userProfile := &UserProfileDTO{
		Account:      user.Account,
		Name:         user.UserProfile.Name,
		Introduction: user.UserProfile.Introduction,
		Avatar:       user.UserProfile.Avatar,
		Role:         user.Role,
		IsMuted:      user.UserProfile.IsMuted,
		Posts:        userPostDTO,
	}
	return userProfile, nil
}

func DeletePost(account string, postID uint, role int) (error, bool) {
	//先判断是否为管理员是否为作者文章
	user, err := global.Post.GetPostDetail(postID)
	if err != nil {
		return err, false
	}
	userAccount := user.User.Account
	if userAccount == account || role == model.RoleAdmin {
		err = global.Post.DeletePost(postID)
		if err != nil {
			return err, false
		}
		return nil, true
	}
	return nil, false
}

func DeleteComment(account string, commentID uint, role int) (error, bool) {
	comment, err := global.Post.GetCommentDetail(commentID)
	if err != nil {
		return err, false
	}
	commentAccount := comment.User.Account
	post, err := global.Post.GetPostDetail(comment.PostID)
	if err != nil {
		return err, false
	}
	posterAccount := post.User.Account
	if commentAccount == account || posterAccount == account || role == model.RoleAdmin {
		err = global.Post.DeleteComment(commentID)
		if err != nil {
			return err, false
		}
		return nil, true
	}
	return nil, false
}
