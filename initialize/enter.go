package initialize

import (
	"skytakeout/config"
	"skytakeout/global"
	"skytakeout/logger"

	"github.com/gin-gonic/gin"
)

func GlobalInit() *gin.Engine {
	// 配置文件初始化
	global.Config = config.InitLoadConfig()
	// Log初始化
	global.Log = logger.NewLogger(global.Config.Log.Level, global.Config.Log.FilePath)

	// Gorm初始化
	global.DB = InitDatabase(global.Config.DataSource.Dsn())
	// Redis初始化
	global.Redis = InitRedis()
	//Router初始化
	router := routerInit()
	return router
}
