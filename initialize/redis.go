package initialize

import (
	"context"
	"fmt"
	"skytakeout/global"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	redisOpt := global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisOpt.Host, redisOpt.Port),
		Password: redisOpt.Password, // no password set
		DB:       redisOpt.DataBase, // use default DB
	})
	ping := client.Ping(context.Background())
	err := ping.Err()
	if err != nil {
		panic(err)
	}
	return client
}
