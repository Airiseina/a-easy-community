package red

import (
	"commmunity/app/internal/model"
	"commmunity/app/utils"
	"commmunity/app/zlog"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (rdb Redis) Like(key string, account string) error {
	err := rdb.redis.SAdd(rdb.context, key, account).Err()
	if err != nil {
		zlog.Error("添加点赞缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) Unlike(key string, account string) error {
	err := rdb.redis.SRem(rdb.context, key, account).Err()
	if err != nil {
		zlog.Error("删除点赞缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) IsLike(key string, account string) (bool, error) {
	isLike, err := rdb.redis.SIsMember(rdb.context, key, account).Result()
	if err != nil {
		zlog.Error("查找是否点赞失败", zap.Error(err))
		return false, err
	}
	return isLike, nil
}

func (rdb Redis) LikeCount(key string) (int64, error) {
	count, err := rdb.redis.SCard(rdb.context, key).Result()
	if err != nil {
		zlog.Error("统计点赞数失败", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (rdb Redis) ScanRedis(match string, cursor uint64) (uint64, []string, error) {
	result := rdb.redis.Scan(rdb.context, cursor, match, 100)
	if result.Err() != nil {
		zlog.Error("查找缓存失败", zap.Error(result.Err()))
		return 0, nil, result.Err()
	}
	key, nextCursor := result.Val()
	return nextCursor, key, nil
}

func (rdb Redis) RateLimiting(ctx context.Context, key string) (int64, error) {
	return rdb.redis.Incr(ctx, key).Result()
}

func (rdb Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rdb.redis.Expire(ctx, key, expiration).Err()
}

func (rdb Redis) LimitView(key string) (bool, error) {
	exist, err := rdb.redis.SetNX(rdb.context, key, 1, 5*time.Minute).Result() //一个账户限制5分钟算一个播放量
	if err != nil {
		zlog.Error("限制播放量出问题", zap.Error(err))
		return false, err
	}
	return exist, nil
}

func (rdb Redis) View(key string) error {
	err := rdb.redis.Incr(rdb.context, key).Err()
	if err != nil {
		zlog.Error("播放量增加出问题", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) ViewCount(key string) (int, error) {
	strCount, err := rdb.redis.Get(rdb.context, key).Result()
	if err != nil {
		zlog.Error("查询播放量出问题", zap.Error(err))
		return 0, err
	}
	count, err := strconv.Atoi(strCount)
	if err != nil {
		zlog.Error("播放量类型转换失败", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (rdb Redis) HotRank(posts []model.Post) error {
	var zMembers []redis.Z
	for _, post := range posts {
		score := float64(post.LikeCount)*1.5 + float64(post.CommentCount)*1 + float64(post.ViewCount)*0.1
		zMembers = append(zMembers, redis.Z{
			Score:  score,
			Member: post.ID,
		})
	}
	err := rdb.redis.ZAdd(rdb.context, "rank:hot", zMembers...).Err()
	if err != nil {
		zlog.Error("存入热度缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) GetHotRank() ([]redis.Z, error) {
	results, err := rdb.redis.ZRevRangeWithScores(rdb.context, "rank:hot", 0, 9).Result()
	if err != nil {
		zlog.Error("提取热度前10失败", zap.Error(err))
		return nil, err
	}
	return results, err
}

func (rdb Redis) SetPostCache(postId uint, postDetail interface{}) error {
	key := fmt.Sprintf("post:cache:%d", postId)
	data, err := json.Marshal(postDetail)
	if err != nil {
		zlog.Error("JSON序列化失败", zap.Error(err))
		return err
	}
	expiration := 1*time.Hour + utils.RandomDuration(5)
	if string(data) == "{}" {
		expiration = 5*time.Minute + utils.RandomDuration(1)
	}
	err = rdb.redis.Set(rdb.context, key, data, expiration).Err()
	if err != nil {
		zlog.Error("建立文章缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) GetPostCache(postId uint) (string, error) {
	key := fmt.Sprintf("post:cache:%d", postId)
	data, err := rdb.redis.Get(rdb.context, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			zlog.Error("获取文章缓存失败", zap.Error(err))
			return "", err
		}
		return "", nil
	}
	return data, nil
}

func (rdb Redis) DelPostCache(postId uint) error {
	key := fmt.Sprintf("post:cache:%d", postId)
	err := rdb.redis.Del(rdb.context, key).Err()
	if err != nil {
		zlog.Error("删除文章缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) GetPostListCache(offset, pageSize int) (string, error) {
	key := fmt.Sprintf("post:list:%d:%d", offset, pageSize)
	data, err := rdb.redis.Get(rdb.context, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			zlog.Error("获取文章列表缓存失败", zap.Error(err))
			return "", err
		}
		return "", nil
	}
	return data, nil
}

func (rdb Redis) SetPostListCache(offset, pageSize int, posts interface{}) error {
	key := fmt.Sprintf("post:list:%d:%d", offset, pageSize)
	data, err := json.Marshal(posts)
	if err != nil {
		zlog.Error("JSON序列化失败", zap.Error(err))
		return err
	}
	err = rdb.redis.Set(rdb.context, key, data, 5*time.Minute+utils.RandomDuration(1)).Err()
	if err != nil {
		zlog.Error("建立文章列表缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) SetFollowingPostsCache(account string, offset, pageSize int, posts interface{}) error {
	key := fmt.Sprintf("following:posts:%s:%d:%d", account, offset, pageSize)
	data, err := json.Marshal(posts)
	if err != nil {
		zlog.Error("JSON序列化失败", zap.Error(err))
		return err
	}
	err = rdb.redis.Set(rdb.context, key, data, 3*time.Minute+utils.RandomDuration(1)).Err()
	if err != nil {
		zlog.Error("建立动态列表缓存失败", zap.Error(err))
		return err
	}
	return nil
}

func (rdb Redis) GetFollowingPostsCache(account string, offset, pageSize int) (string, error) {
	key := fmt.Sprintf("following:posts:%s:%d:%d", account, offset, pageSize)
	data, err := rdb.redis.Get(rdb.context, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			zlog.Error("获取动态列表缓存失败", zap.Error(err))
			return "", err
		}
		return "", nil
	}
	return data, nil
}
