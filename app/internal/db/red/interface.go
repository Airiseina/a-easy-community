package red

import "time"

type UserRedis interface {
	AddToBlacklist(tokenString string, expiration time.Duration) error
	IsInBlacklist(tokenString string) bool
}
