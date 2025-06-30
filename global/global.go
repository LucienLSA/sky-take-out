package global

import (
	"skytakeout/config"

	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"gorm.io/gorm"
)

var (
	Config *config.AllConfig
	// Log    logger.ILog
	DB     *gorm.DB
	Rdb    *redis.Client
	ZapLog *otelzap.Logger
)
