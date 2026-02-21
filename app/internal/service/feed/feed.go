package feed

import (
	"commmunity/app/internal/db/global"
	"encoding/json"
)

func Follow(account string, followedId uint) (error, bool) {
	follower, err := global.User.GetUserId(account)
	if err != nil {
		return err, false
	}
	followerId := follower.ID
	flag, err := global.User.IsFollowing(followedId, followerId)
	if err != nil {
		return err, false
	}
	err = global.UserRedis.DelFollowersCache(account)
	if err != nil {
		return err, false
	}
	err = global.UserRedis.DelFollowingsCache(account)
	if err != nil {
		return err, false
	}
	if flag {
		err = global.User.Unfollow(followedId, followerId)
		if err != nil {
			return err, false
		}
		return nil, false
	} else {
		err = global.User.Follow(followedId, followerId)
		if err != nil {
			return err, false
		}
		return nil, true
	}
}

type FollowsDTO struct {
	UserId       uint   `json:"user_id"`
	Avatar       string `json:"avatar"`
	Name         string `json:"name"`
	Introduction string `json:"introduction"`
}

func GetFollowers(account string) ([]FollowsDTO, error) {
	fc, err := global.UserRedis.GetFollowersCache(account)
	if err != nil {
		return nil, err
	}
	if fc == "[]" {
		return []FollowsDTO{}, nil
	}
	if fc != "" {
		var cached []FollowsDTO
		if err = json.Unmarshal([]byte(fc), &cached); err == nil {
			return cached, nil
		}
	}
	var followersDTO []FollowsDTO
	user, err := global.User.GetUserId(account)
	if err != nil {
		return nil, err
	}
	followers, err := global.User.GetFollowers(user.ID)
	if err != nil {
		return nil, err
	}
	for _, follower := range followers {
		followerDTO := FollowsDTO{
			UserId:       follower.ID,
			Avatar:       follower.UserProfile.Avatar,
			Name:         follower.UserProfile.Name,
			Introduction: follower.UserProfile.Introduction,
		}
		followersDTO = append(followersDTO, followerDTO)
	}
	if len(followersDTO) == 0 {
		err = global.UserRedis.SetFollowersCache(account, []FollowsDTO{})
		if err != nil {
			return nil, err
		}
	} else {
		err = global.UserRedis.SetFollowersCache(account, followersDTO)
		if err != nil {
			return nil, err
		}
	}
	return followersDTO, nil
}

func GetFollowings(account string) ([]FollowsDTO, error) {
	fc, err := global.UserRedis.GetFollowingsCache(account)
	if err != nil {
		return nil, err
	}
	if fc == "[]" {
		return []FollowsDTO{}, nil
	}
	if fc != "" {
		var cached []FollowsDTO
		if err = json.Unmarshal([]byte(fc), &cached); err == nil {
			return cached, nil
		}
	}
	var followingsDTO []FollowsDTO
	user, err := global.User.GetUserId(account)
	if err != nil {
		return nil, err
	}
	followings, err := global.User.GetFollowings(user.ID)
	if err != nil {
		return nil, err
	}
	for _, following := range followings {
		followingDTO := FollowsDTO{
			UserId:       following.ID,
			Avatar:       following.UserProfile.Avatar,
			Name:         following.UserProfile.Name,
			Introduction: following.UserProfile.Introduction,
		}
		followingsDTO = append(followingsDTO, followingDTO)
	}

	if len(followingsDTO) == 0 {
		err = global.UserRedis.SetFollowingsCache(account, []FollowsDTO{})
		if err != nil {
			return nil, err
		}
	} else {
		err = global.UserRedis.SetFollowingsCache(account, followingsDTO)
		if err != nil {
			return nil, err
		}
	}

	return followingsDTO, nil
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

func GetFollowingPosts(account string, offset int, pageSize int) ([]PostsDTO, error) {
	fpc, err := global.PostRedis.GetFollowingPostsCache(account, offset, pageSize)
	if err != nil {
		return nil, err
	}
	if fpc == "[]" {
		return []PostsDTO{}, nil
	}
	if fpc != "" {
		var cached []PostsDTO
		if err = json.Unmarshal([]byte(fpc), &cached); err == nil {
			return cached, nil
		}
	}
	user, err := global.User.GetUserId(account)
	if err != nil {
		return nil, err
	}
	ps, err := global.Post.GetFollowingPosts(user.ID, offset, pageSize)
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
	if len(posts) == 0 {
		err = global.PostRedis.SetFollowingPostsCache(account, offset, pageSize, []PostsDTO{})
		if err != nil {
			return nil, err
		}
	} else {
		err = global.PostRedis.SetFollowingPostsCache(account, offset, pageSize, posts)
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}
