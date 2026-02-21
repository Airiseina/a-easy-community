package red

import (
	"commmunity/app/utils"
	"commmunity/app/zlog"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (rdb Redis) AddToBlacklist(tokenString string, expiration time.Duration) error {
	return rdb.redis.Set(rdb.context, "blacklist:"+tokenString, "1", expiration).Err()
}

func (rdb Redis) IsInBlacklist(tokenString string) bool {
	n, _ := rdb.redis.Exists(rdb.context, "blacklist:"+tokenString).Result()
	return n > 0
}

func (rdb Redis) UserProfile(userId uint, userDetail interface{}) error {
	key := fmt.Sprintf("user:cache:%d", userId)
	data, err := json.Marshal(userDetail)
	if err != nil {
		zlog.Error("JSON序列化失败", zap.Error(err))
		return err
	}
	expiration := 24*time.Hour + utils.RandomDuration(5)
	if string(data) == "{}" {
		expiration = 5*time.Minute + utils.RandomDuration(1)
	}
	err = rdb.redis.Set(rdb.context, key, data, expiration).Err()
	if err != nil {
		zlog.Error("建立用户信息缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) GetUserCache(userId uint) (string, error) {
	key := fmt.Sprintf("user:cache:%d", userId)
	data, err := rdb.redis.Get(rdb.context, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			zlog.Error("获取用户信息缓存失败", zap.Error(err))
			return "", err
		}
		return "", nil
	}
	return data, nil
}

func (rdb Redis) DelUserCache(userId uint) error {
	key := fmt.Sprintf("user:cache:%d", userId)
	err := rdb.redis.Del(rdb.context, key).Err()
	if err != nil {
		zlog.Error("删除用户缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) SetFollowersCache(account string, followers interface{}) error {
	key := fmt.Sprintf("followers:%s", account)
	data, err := json.Marshal(followers)
	if err != nil {
		zlog.Error("JSON序列化失败", zap.Error(err))
		return err
	}
	err = rdb.redis.Set(rdb.context, key, data, 10*time.Minute+utils.RandomDuration(2)).Err()
	if err != nil {
		zlog.Error("建立粉丝列表缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) GetFollowersCache(account string) (string, error) {
	key := fmt.Sprintf("followers:%s", account)
	data, err := rdb.redis.Get(rdb.context, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			zlog.Error("获取粉丝列表缓存失败", zap.Error(err))
			return "", err
		}
		return "", nil
	}
	return data, nil
}

func (rdb Redis) DelFollowersCache(account string) error {
	key := fmt.Sprintf("followers:%s", account)
	err := rdb.redis.Del(rdb.context, key).Err()
	if err != nil {
		zlog.Error("删除粉丝列表缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) SetFollowingsCache(account string, followings interface{}) error {
	key := fmt.Sprintf("followings:%s", account)
	data, err := json.Marshal(followings)
	if err != nil {
		zlog.Error("JSON序列化失败", zap.Error(err))
		return err
	}
	err = rdb.redis.Set(rdb.context, key, data, 10*time.Minute+utils.RandomDuration(2)).Err()
	if err != nil {
		zlog.Error("建立关注列表缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) GetFollowingsCache(account string) (string, error) {
	key := fmt.Sprintf("followings:%s", account)
	data, err := rdb.redis.Get(rdb.context, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			zlog.Error("获取关注列表缓存失败", zap.Error(err))
			return "", err
		}
		return "", nil
	}
	return data, nil
}

func (rdb Redis) DelFollowingsCache(account string) error {
	key := fmt.Sprintf("followings:%s", account)
	err := rdb.redis.Del(rdb.context, key).Err()
	if err != nil {
		zlog.Error("删除关注列表缓存失败", zap.Error(err))
		return err
	}
	return nil
}
