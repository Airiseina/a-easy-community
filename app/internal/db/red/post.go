package red

import (
	"commmunity/app/zlog"

	"go.uber.org/zap"
)

func (rdb Redis) Like(key string, userId uint) error {
	err := rdb.redis.SAdd(rdb.context, key, userId).Err()
	if err != nil {
		zlog.Error("添加点赞缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) Unlike(key string, userId uint) error {
	err := rdb.redis.SRem(rdb.context, key, userId).Err()
	if err != nil {
		zlog.Error("删除点赞缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) IsLike(key string, userId uint) (bool, error) {
	isLike, err := rdb.redis.SIsMember(rdb.context, key, userId).Result()
	if err != nil {
		zlog.Error("查找是否点赞失败", zap.Error(err))
		return false, err
	}
	return isLike, nil
}
