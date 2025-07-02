package cache

import (
	"context"
	"fmt"
	"skytakeout/common/e"
	"skytakeout/common/retcode"
	"skytakeout/global"
	"time"

	"skytakeout/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 设置用户会话失效通知
func SetSessionInvalidNotification(ctx context.Context, username string) error {
	notificationKey := fmt.Sprintf("session:invalid:%s", username)
	// 设置通知，TTL为5分钟，足够客户端轮询检查
	err := global.Rdb.Set(ctx, notificationKey, "invalid", 5*time.Minute).Err()
	if err != nil {
		logger.Logger(ctx).Error("global.Rdb.Set failed", zap.Error(err))
		return retcode.NewError(e.RedisERR, "rdb.Set failed")
	}
	logger.Logger(ctx).Info("设置会话失效通知", zap.String("username", username))
	return nil
}

// 检查用户是否有会话失效通知
func HasSessionInvalidNotification(ctx context.Context, username string) (bool, error) {
	notificationKey := fmt.Sprintf("session:invalid:%s", username)
	_, err := global.Rdb.Get(ctx, notificationKey).Result()
	if err == redis.Nil {
		return false, nil // 没有通知
	}
	if err != nil {
		logger.Logger(ctx).Error("global.Rdb.Get failed", zap.Error(err))
		return false, retcode.NewError(e.RedisERR, "rdb.Get failed")
	}
	return true, nil // 有通知
}

// 清除会话失效通知
func ClearSessionInvalidNotification(ctx context.Context, username string) error {
	notificationKey := fmt.Sprintf("session:invalid:%s", username)
	err := global.Rdb.Del(ctx, notificationKey).Err()
	if err != nil {
		logger.Logger(ctx).Error("global.Rdb.Del failed", zap.Error(err))
		return retcode.NewError(e.RedisERR, "rdb.Del failed")
	}
	return nil
}
