package red

import (
	"time"
)

func (rdb Redis) AddToBlacklist(tokenString string, expiration time.Duration) error {
	return rdb.redis.Set(rdb.context, "blacklist:"+tokenString, "1", expiration).Err()
}

func (rdb Redis) IsInBlacklist(tokenString string) bool {
	n, _ := rdb.redis.Exists(rdb.context, "blacklist:"+tokenString).Result()
	return n > 0
}
