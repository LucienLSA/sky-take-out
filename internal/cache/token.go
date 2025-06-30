package cache

import (
	"context"
	"fmt"
	"skytakeout/common/e"
	"skytakeout/common/retcode"
	"skytakeout/global"
	"time"

	"github.com/redis/go-redis/v9"
)

func StoreUserIdToken(ctx context.Context, token, username string) (err error) {
	// 获取token的存活时间
	duration := time.Duration(24 * time.Hour)
	key := GetRedisKey(KeyTokenSetPrefix)
	// 存入redis
	if err = global.Rdb.Set(ctx, key+username, token, duration).Err(); err != nil {
		global.Log.Error(ctx, "global.Rdb.Set failed, err: %v", err)
		return retcode.NewError(e.RedisERR, "rdb.Set failed")
	}
	return nil
}

// 从redis取token
func GetJwtToken(ctx context.Context, username string) (token string, err error) {
	key := GetRedisKey(KeyTokenSetPrefix)
	token, err = global.Rdb.Get(ctx, key+username).Result()
	fmt.Println("token:", token)
	if err == redis.Nil {
		global.Log.Error(ctx, "global.Rdb.Get failed, err: redis.Nil")
		return "", retcode.NewError(e.ErrorUserNotLogin, "token is empty")
	}
	if err != nil {
		global.Log.Error(ctx, "global.Rdb.Get failed, err: %v", err)
		return "", retcode.NewError(e.RedisERR, "rdb.Get failed")
	}
	return token, nil
}

// 删除用户token
func DeleteUserIdToken(ctx context.Context, username string) error {
	// 由于存储时是以 username 为 key，token 为 value，需要查找对应username的token进行过期
	// username在数据库定义时是唯一值
	// 构造Redis中存储的key
	key := GetRedisKey(KeyTokenSetPrefix) + username
	// 直接删除该key todo:设置过期时间为0
	err := global.Rdb.Del(ctx, key).Err()
	if err != nil {
		global.Log.Error(ctx, "global.Rdb.De, err: %v", err)
		return retcode.NewError(e.RedisERR, "rdb.Del failed")
	}
	return nil
}
