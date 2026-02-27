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

func (rdb Redis) SetMessageCache(userId1 uint, userId2 uint, value interface{}, offset, pageSize int) error {
	key := fmt.Sprintf("message_cache:%d:%d:%d:%d", userId1, userId2, offset, pageSize)
	key1 := fmt.Sprintf("message_cache:%d:%d:%d:%d", userId2, userId1, offset, pageSize)
	data, err := json.Marshal(value)
	if err != nil {
		zlog.Error("JSON序列化失败", zap.Error(err))
		return err
	}
	err = rdb.redis.Set(rdb.context, key, data, 30*time.Minute+utils.RandomDuration(1)).Err()
	if err != nil {
		zlog.Error("建立历史消息缓存失败", zap.Error(err))
		return err
	}
	err = rdb.redis.Set(rdb.context, key1, data, 30*time.Minute+utils.RandomDuration(1)).Err()
	if err != nil {
		zlog.Error("建立历史消息缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) GetMessageCache(userId1 uint, userId2 uint, offset int, pageSize int) (string, error) {
	key := fmt.Sprintf("message_cache:%d:%d:%d:%d", userId1, userId2, offset, pageSize)
	data, err := rdb.redis.Get(rdb.context, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			zlog.Error("获取历史消息缓存失败", zap.Error(err))
			return "", err
		}
		return "", nil
	}
	return data, nil
}

func (rdb Redis) DelMessageCache(userId1 uint, userId2 uint) error {
	pattern := fmt.Sprintf("message_cache:%d:%d:*", userId1, userId2)
	iter := rdb.redis.Scan(rdb.context, 0, pattern, 0).Iterator()
	if err := iter.Err(); err != nil {
		zlog.Error("遍历缓存key失败", zap.Error(err))
		return err
	}
	for iter.Next(rdb.context) {
		err := rdb.redis.Del(rdb.context, iter.Val()).Err()
		if err != nil {
			zlog.Error("删除缓存失败", zap.Error(err), zap.String("key", iter.Val()))
		}
	}
	pattern2 := fmt.Sprintf("message_cache:%d:%d:*", userId2, userId1)
	iter2 := rdb.redis.Scan(rdb.context, 0, pattern2, 0).Iterator()
	if err := iter2.Err(); err != nil {
		zlog.Error("遍历缓存key失败", zap.Error(err))
		return err
	}
	for iter2.Next(rdb.context) {
		err := rdb.redis.Del(rdb.context, iter2.Val()).Err()
		if err != nil {
			zlog.Error("删除缓存失败", zap.Error(err), zap.String("key", iter2.Val()))
		}
	}
	return nil
}
