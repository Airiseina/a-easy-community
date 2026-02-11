package feed

import "commmunity/app/internal/db/global"

func Follow(account string, followedId uint) (error, bool) {
	followed, err := global.User.GetUserId(account)
	if err != nil {
		return err, false
	}
	followerId := followed.ID
	flag, err := global.User.IsFollowing(followedId, followerId)
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
	return followersDTO, nil
}

func GetFollowings(account string) ([]FollowsDTO, error) {
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
	return posts, nil
}
