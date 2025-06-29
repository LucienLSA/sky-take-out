package main

import (
	"skytakeout/global"
	"skytakeout/initialize"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	router := initialize.GlobalInit()

	// 设置运行环境
	gin.SetMode(global.Config.Server.Level)

	router.Run(":8080")
}
