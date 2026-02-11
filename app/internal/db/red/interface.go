package red

import (
	"context"
	"time"
)

type UserRedis interface {
	AddToBlacklist(tokenString string, expiration time.Duration) error
	IsInBlacklist(tokenString string) bool
}

type PostRedis interface {
	Like(key string, account string) error
	Unlike(key string, account string) error
	IsLike(key string, account string) (bool, error)
	LikeCount(key string) (int64, error)
	ScanLikes(match string, cursor uint64) (uint64, []string, error)
	RateLimiting(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	View(key string) error
	ViewCount(key string) (int, error)
}
