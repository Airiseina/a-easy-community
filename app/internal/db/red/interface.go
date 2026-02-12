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
}
