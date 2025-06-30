package main

import (
	"context"
	"skytakeout/config"
	"skytakeout/global"
	"skytakeout/initialize"
	"skytakeout/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 配置文件初始化
	global.Config = config.InitLoadConfig()
	// Log初始化
	// global.Log = logger.NewLogger(global.Config.Log.Level, global.Config.Log.FilePath)
	logger.Init(global.Config.Log.Level, global.Config.Log.FilePath)
	logger.Logger(context.Background()).Info("logger init success")
	// tracer 初始化
	shutdown := initialize.InitTracer()
	defer shutdown(context.Background())
	logger.Logger(context.Background()).Info("trace init success")
	// Gorm初始化
	global.DB = initialize.InitDatabase(global.Config.DataSource.Dsn())
	logger.Logger(context.Background()).Info("gorm init success")
	// Redis初始化
	global.Rdb = initialize.InitRedis()
	logger.Logger(context.Background()).Info("redis init success")
	// routerRedis初始化
	router := initialize.InitRouter()
	logger.Logger(context.Background()).Info("route init success")
	// 设置运行环境
	gin.SetMode(global.Config.Server.Level)

	router.Run(":8080")
}
