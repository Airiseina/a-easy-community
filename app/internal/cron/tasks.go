package cron

import (
	"commmunity/app/internal/db/global"
	"commmunity/app/zlog"
	"context"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func SyncPostLikes(ctx context.Context) {
	key := "post:likes:*"
	cursor := uint64(0)
	var uniqueKey []string
	seen := make(map[string]bool)
	for {
		select {
		case <-ctx.Done():
			zlog.Info("点赞同步任务被取消，停止扫描")
			return
		default:
		}
		nextCursor, appendKeys, err := global.PostRedis.ScanLikes(key, cursor)
		if err != nil {
			zlog.Error("扫描点赞失败", zap.Error(err))
			break
		}
		for _, appendKey := range appendKeys {
			if !seen[appendKey] {
				seen[appendKey] = true
				uniqueKey = append(uniqueKey, appendKey)
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	for _, k := range uniqueKey {
		select {
		case <-ctx.Done():
			zlog.Info("点赞同步任务被取消，停止更新")
			return
		default:
		}
		parts := strings.Split(k, ":")
		if len(parts) < 3 {
			continue
		}
		postIdStr := parts[2]
		postIdInt, err := strconv.ParseUint(postIdStr, 10, 64)
		if err != nil {
			zlog.Error("解析PostId失败", zap.Error(err))
			continue
		}
		postId := uint(postIdInt)
		c, err := global.PostRedis.LikeCount(k)
		if err != nil {
			zlog.Error("获取Redis点赞数失败", zap.Error(err))
			continue
		}
		count := uint(c)
		err = global.Post.Like(postId, count)
		if err != nil {
			zlog.Error("同步Post点赞数失败", zap.Error(err))
			continue
		}
	}
}
