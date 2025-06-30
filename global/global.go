package global

import (
	"skytakeout/config"
	"skytakeout/logger"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config *config.AllConfig
	Log    logger.ILog
	DB     *gorm.DB
	Rdb    *redis.Client
)
