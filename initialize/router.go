package initialize

import (
	"skytakeout/global"
	"skytakeout/internal/router"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	allRouter := router.AllRouter

	// 链路追踪日志中间件
	r.Use(otelgin.Middleware(global.Config.Jaeger.ServiceName))

	// admin
	admin := r.Group("/admin")
	{
		allRouter.EmployeeRouter.InitApiRouter(admin)
		allRouter.CategoryRouter.InitApiRouter(admin)
		allRouter.DishRouter.InitApiRouter(admin)
		allRouter.CommonRouter.InitApiRouter(admin)
		allRouter.SetMealRouter.InitApiRouter(admin)
	}
	return r
}
