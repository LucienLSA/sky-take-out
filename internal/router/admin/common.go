package admin

import (
	"github.com/gin-gonic/gin"
)

type CommonRouter struct{}

func (dr *CommonRouter) InitApiRouter(parent *gin.RouterGroup) {
	//publicRouter := parent.Group("category")
	// privateRouter := parent.Group("common")
	// 私有路由使用jwt验证
	// privateRouter.Use(middle.VerifyJWTAdmin())
	// 依赖注入
	// commonCtrl := new(controller.CommonController)
	{
		// privateRouter.POST("upload", commonCtrl.Upload)
	}
}
