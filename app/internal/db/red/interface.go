package red

import (
	"commmunity/app/internal/model"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserRedis interface {
	AddToBlacklist(tokenString string, expiration time.Duration) error
	IsInBlacklist(tokenString string) bool
	UserProfile(userId uint, userDetail interface{}) error
	GetUserCache(userId uint) (string, error)
	DelUserCache(userId uint) error
	SetFollowersCache(account string, followers interface{}) error
	GetFollowersCache(account string) (string, error)
	DelFollowersCache(account string) error
	SetFollowingsCache(account string, followings interface{}) error
	GetFollowingsCache(account string) (string, error)
	DelFollowingsCache(account string) error
}

type PostRedis interface {
	Like(key string, account string) error
	Unlike(key string, account string) error
	IsLike(key string, account string) (bool, error)
	LikeCount(key string) (int64, error)
	ScanRedis(match string, cursor uint64) (uint64, []string, error)
	RateLimiting(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	LimitView(key string) (bool, error)
	View(key string) error
	ViewCount(key string) (int, error)
	HotRank(posts []model.Post) error
	GetHotRank() ([]redis.Z, error)
	SetPostCache(postId uint, postDetail interface{}) error
	GetPostCache(postId uint) (string, error)
	DelPostCache(postId uint) error
	SetPostListCache(offset, pageSize int, posts interface{}) error
	GetPostListCache(offset, pageSize int) (string, error)
	SetFollowingPostsCache(account string, offset, pageSize int, posts interface{}) error
	GetFollowingPostsCache(account string, offset, pageSize int) (string, error)
	SetSummaryCache(postId uint, summary string) error
	GetSummaryCache(postId uint) (string, error)
}
