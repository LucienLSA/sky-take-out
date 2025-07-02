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

func StoreUserAToken(ctx context.Context, token, username string) (err error) {
	accessKey := fmt.Sprintf("jwt:admin:%s:access", username)
	ttl := global.Config.Jwt.Admin.TTL
	accessDuration := time.Duration(ttl) * time.Minute
	// 存入redis
	if err = global.Rdb.Set(ctx, accessKey, token, accessDuration).Err(); err != nil {
		logger.Logger(ctx).Error("global.Rdb.Set failed", zap.Error(err))
		return retcode.NewError(e.RedisERR, "rdb.Set failed")
	}
	return nil
}
func StoreUserRToken(ctx context.Context, token, username string) (err error) {
	ttl := global.Config.Jwt.Admin.TTL
	refreshDuration := time.Duration(ttl) * time.Hour
	refreshKey := fmt.Sprintf("jwt:admin:%s:refresh", username)
	// 存入redis
	if err = global.Rdb.Set(ctx, refreshKey, token, refreshDuration).Err(); err != nil {
		logger.Logger(ctx).Error("global.Rdb.Set failed", zap.Error(err))
		return retcode.NewError(e.RedisERR, "rdb.Set failed")
	}
	return nil
}

// 删除用户access_token
func DeleteUserAToken(ctx context.Context, username string) error {
	accessKey := fmt.Sprintf("jwt:admin:%s:access", username)
	err := global.Rdb.Del(ctx, accessKey).Err()
	if err != nil {
		logger.Logger(ctx).Error("global.Rdb.Del failed", zap.Error(err))
		return retcode.NewError(e.RedisERR, "rdb.Del failed")
	}
	return nil
}

// 删除用户refresh_token
func DeleteUserRToken(ctx context.Context, username string) error {
	refreshKey := fmt.Sprintf("jwt:admin:%s:refresh", username)
	err := global.Rdb.Del(ctx, refreshKey).Err()
	if err != nil {
		logger.Logger(ctx).Error("global.Rdb.Del failed", zap.Error(err))
		return retcode.NewError(e.RedisERR, "rdb.Del failed")
	}
	return nil
}

// 从redis取access_token
func GetUserAToken(ctx context.Context, username string) (token string, err error) {
	accessKey := fmt.Sprintf("jwt:admin:%s:access", username)
	token, err = global.Rdb.Get(ctx, accessKey).Result()
	if err == redis.Nil {
		logger.Logger(ctx).Error("global.Rdb.Get failed", zap.String("err", "redis.Nil"))
		return "", retcode.NewError(e.ErrorUserNotLogin, "token is empty")
	}
	if err != nil {
		logger.Logger(ctx).Error("global.Rdb.Get failed", zap.Error(err))
		return "", retcode.NewError(e.RedisERR, "rdb.Get failed")
	}
	return token, nil
}

// 强制清除用户所有会话（用于单点登录）
func ForceLogoutUser(ctx context.Context, username string) error {
	// 删除access_token
	err := DeleteUserAToken(ctx, username)
	if err != nil {
		logger.Logger(ctx).Error("DeleteUserAToken failed", zap.Error(err))
		return err
	}

	// 删除refresh_token
	err = DeleteUserRToken(ctx, username)
	if err != nil {
		logger.Logger(ctx).Error("DeleteUserRToken failed", zap.Error(err))
		return err
	}

	logger.Logger(ctx).Info("用户会话已强制清除", zap.String("username", username))
	return nil
}
