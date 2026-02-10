package red

import "time"

type UserRedis interface {
	AddToBlacklist(tokenString string, expiration time.Duration) error
	IsInBlacklist(tokenString string) bool
}

type PostRedis interface {
	Like(key string, userId uint) error
}
