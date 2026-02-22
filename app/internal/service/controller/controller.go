package controller

import (
	"commmunity/app/internal/ai"
	"commmunity/app/internal/db/global"
	"commmunity/app/internal/model"
	"commmunity/app/utils"
	"commmunity/app/zlog"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var requestGroup singleflight.Group

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
	Paid         bool   `json:"paid"`
	ViewCount    uint   `json:"view_count"`
	LikeCount    uint   `json:"like_count"`
	CommentCount uint   `json:"comment_count"`
}

func GetPostList(offset int, pageSize int) ([]PostsDTO, error) {
	pc, err := global.PostRedis.GetPostListCache(offset, pageSize)
	if err != nil {
		return nil, err
	}
	if pc == "[]" {
		return []PostsDTO{}, nil
	}
	if pc != "" {
		var cachedPosts []PostsDTO
		if err = json.Unmarshal([]byte(pc), &cachedPosts); err == nil {
			return cachedPosts, nil
		} else {
			return nil, err
		}
	}
	val, err, _ := requestGroup.Do(fmt.Sprintf("post:list:%d:%d", offset, pageSize), func() (interface{}, error) {
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
				Paid:         p.Paid,
				ViewCount:    p.ViewCount,
				LikeCount:    p.LikeCount,
				CommentCount: p.CommentCount,
			}
		}
		if len(posts) == 0 {
			err = global.PostRedis.SetPostListCache(offset, pageSize, []PostsDTO{})
			if err != nil {
				return nil, err
			}
		}
		err = global.PostRedis.SetPostListCache(offset, pageSize, posts)
		if err != nil {
			return nil, err
		}
		return posts, nil
	})
	if err != nil {
		return nil, err
	}
	return val.([]PostsDTO), nil
}

type PostDTO struct {
	PostsDTO
	Content   string       `json:"content"`
	UserId    uint         `json:"user_id"`
	Comment   []CommentDTO `json:"comment"`
	AiSummary string       `json:"ai_summary"`
}

type CommentDTO struct {
	ID         uint   `json:"id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
	UserId     uint   `json:"user_id"`
}

func GetPostDetail(account string, postId uint) (PostDTO, error) {
	pc, err := global.PostRedis.GetPostCache(postId)
	if err != nil {
		return PostDTO{}, err
	}
	if pc == "{}" {
		return PostDTO{}, nil
	}
	if pc != "" {
		var cachedPost PostDTO
		if err = json.Unmarshal([]byte(pc), &cachedPost); err == nil {
			key := fmt.Sprintf("post:view:%d", postId)
			limitKey := fmt.Sprintf("post:view:limit:%s:%d", account, postId)
			flag, err := global.PostRedis.LimitView(limitKey)
			if err != nil {
				return PostDTO{}, err
			}
			if flag {
				err = global.PostRedis.View(key)
				if err != nil {
					return PostDTO{}, err
				}
			}
			viewCount, err := global.PostRedis.ViewCount(key)
			if err != nil {
				return PostDTO{}, err
			}
			likeKey := fmt.Sprintf("post:likes:%d", postId)
			likeCount, err := global.PostRedis.LikeCount(likeKey)
			if err != nil {
				return PostDTO{}, err
			}
			cachedPost.ViewCount = uint(viewCount)
			cachedPost.LikeCount = uint(likeCount)
			return cachedPost, nil
		} else {
			return PostDTO{}, err
		}
	}
	val, err, _ := requestGroup.Do(fmt.Sprintf("post:%d", postId), func() (interface{}, error) {
		user, err := global.User.GetUserId(account)
		p, err := global.Post.GetPostDetail(postId)
		if err != nil {
			return PostDTO{}, err
		}
		if p.ID == 0 {
			_ = global.PostRedis.SetPostCache(postId, map[string]interface{}{})
			return PostDTO{}, nil
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
		postCache := PostDTO{
			PostsDTO: PostsDTO{
				Name:         p.User.UserProfile.Name,
				Avatar:       p.User.UserProfile.Avatar,
				PostID:       p.ID,
				Title:        p.Title,
				Paid:         p.Paid,
				ViewCount:    p.ViewCount,
				LikeCount:    p.LikeCount,
				CommentCount: p.CommentCount,
			},
			Content: p.Content,
			UserId:  p.User.ID,
			Comment: commentDTOs,
		}
		err = global.PostRedis.SetPostCache(postId, postCache)
		if err != nil {
			return PostDTO{}, err
		}
		if !p.Paid || user.Vip {
			post := postCache
			return post, nil
		}
		content := utils.TruncateContent(p.Content, 4, 200)
		post := postCache
		post.Content = content
		return post, nil
	})
	if err != nil {
		return PostDTO{}, err
	}
	post := val.(PostDTO)
	key := fmt.Sprintf("post:view:%d", postId)
	limitKey := fmt.Sprintf("post:view:limit:%s:%d", account, postId)
	flag, err := global.PostRedis.LimitView(limitKey)
	if err != nil {
		return PostDTO{}, err
	}
	if flag {
		err = global.PostRedis.View(key)
		if err != nil {
			return PostDTO{}, err
		}
	}
	return post, nil
}

func AiSummary(account string, postId uint) (string, error) {
	summary, err := global.PostRedis.GetSummaryCache(postId)
	if err != nil {
		return "", err
	}
	if summary != "" {
		return summary, nil
	}
	val, err, _ := requestGroup.Do(fmt.Sprintf("post:%d", postId), func() (interface{}, error) {
		user, err := global.User.GetUserId(account)
		p, err := global.Post.GetPostDetail(postId)
		if err != nil {
			return "", err
		}
		if user.Vip {
			content, err := ai.AutoSummary(p.Content)
			if err != nil {
				return "", err
			}
			err = global.PostRedis.SetSummaryCache(postId, content)
			if err != nil {
				return "", err
			}
			return content, nil
		}
		return "", nil
	})
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func CreateComment(account string, postID uint, content string) (error, bool) {
	user, err := global.User.GetUserId(account)
	if err != nil {
		return err, false
	}
	if user.UserProfile.IsMuted {
		return nil, false
	}
	cleanContent := bluemonday.UGCPolicy().Sanitize(content)
	if err = global.Post.CreateComment(user.ID, postID, cleanContent); err != nil {
		return err, false
	}
	return global.PostRedis.DelPostCache(postID), true
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

func GetUserProfile(id uint) (UserProfileDTO, error) {
	uc, err := global.UserRedis.GetUserCache(id)
	if err != nil {
		return UserProfileDTO{}, err
	}
	if uc == "{}" {
		return UserProfileDTO{}, nil
	}
	if uc != "" {
		var userCache UserProfileDTO
		if err := json.Unmarshal([]byte(uc), &userCache); err == nil {
			return userCache, nil
		} else {
			return UserProfileDTO{}, err
		}
	}
	val, err, _ := requestGroup.Do(fmt.Sprintf("user:%d", id), func() (interface{}, error) {
		user, err := global.Post.GetUserProfile(id)
		if err != nil {
			return UserProfileDTO{}, err
		}
		if user.ID == 0 {
			err = global.UserRedis.UserProfile(id, map[string]interface{}{})
			if err != nil {
				return UserProfileDTO{}, err
			}
			return UserProfileDTO{}, nil
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
		userProfile := UserProfileDTO{
			Account:      user.Account,
			Name:         user.UserProfile.Name,
			Introduction: user.UserProfile.Introduction,
			Avatar:       user.UserProfile.Avatar,
			Role:         user.Role,
			IsMuted:      user.UserProfile.IsMuted,
			Posts:        userPostDTO,
		}
		err = global.UserRedis.UserProfile(id, userProfile)
		if err != nil {
			return UserProfileDTO{}, err
		}
		return userProfile, nil
	})
	if err != nil {
		return UserProfileDTO{}, err
	}
	return val.(UserProfileDTO), nil
}

func DeletePost(account string, postID uint, role int) (error, bool) {
	//先判断是否为管理员，是否为作者文章
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
		err = global.PostRedis.DelPostCache(postID)
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
		err = global.PostRedis.DelPostCache(comment.PostID)
		if err != nil {
			return err, false
		}
		return nil, true
	}
	return nil, false
}

func ToggleLike(postId uint, account string) (bool, int, error) {
	key := fmt.Sprintf("post:likes:%d", postId)
	isLike, err := global.PostRedis.IsLike(key, account)
	if err != nil {
		return false, 0, err
	}
	if isLike {
		err = global.PostRedis.Unlike(key, account)
		if err != nil {
			return false, 0, err
		}
		isLike = false
	} else {
		err = global.PostRedis.Like(key, account)
		if err != nil {
			return false, 0, err
		}
		isLike = true
	}
	c, err := global.PostRedis.LikeCount(key)
	if err != nil {
		return false, 0, err
	}
	count := int(c)
	return isLike, count, nil
}

func RateLimiting(ctx context.Context, key string, limitDuration time.Duration, limitCount int) (bool, error) {
	count, err := global.PostRedis.RateLimiting(ctx, key)
	if err != nil {
		return false, err
	}
	if count == 1 {
		err = global.PostRedis.Expire(ctx, key, limitDuration)
		if err != nil {
			return false, err
		}
	}
	if int(count) > limitCount {
		return false, nil
	}
	return true, nil
}

func GetHotRank() ([]PostsDTO, map[uint]float64, error) {
	var postIds []uint
	score := make(map[uint]float64)
	postsZ, err := global.PostRedis.GetHotRank()
	if err != nil {
		return nil, nil, err
	}
	for _, z := range postsZ {
		idStr := z.Member.(string)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			zlog.Error("解析字符串失败", zap.Error(err))
			return nil, nil, err
		}
		postId := uint(id)
		postIds = append(postIds, postId)
		score[postId] = z.Score
	}
	val, err, _ := requestGroup.Do(fmt.Sprintf("post:hotList"), func() (interface{}, error) {
		ps, err := global.Post.HotPosts(postIds)
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
	})
	return val.([]PostsDTO), score, nil
}

func SearchPosts(keyword string, offset int, pageSize int) ([]PostsDTO, error) {
	posts, err := global.Post.SearchPosts(keyword, offset, pageSize)
	if err != nil {
		return nil, err
	}
	var results []PostsDTO
	for _, p := range posts {
		results = append(results, PostsDTO{
			Name:         p.User.UserProfile.Name,
			Avatar:       p.User.UserProfile.Avatar,
			PostID:       p.ID,
			Title:        p.Title,
			ViewCount:    p.ViewCount,
			LikeCount:    p.LikeCount,
			CommentCount: p.CommentCount,
		})
	}
	return results, nil
}

func SetPostPaid(role int, postId uint) (bool, error) {
	if role == model.RoleAdmin {
		return true, global.Post.SetPostPaid(postId, true)
	}
	return false, nil
}

func PayVip(role int, userId uint) (bool, error) {
	if role == model.RoleAdmin {
		return true, global.User.SetVip(userId, true)
	}
	return false, nil
}
