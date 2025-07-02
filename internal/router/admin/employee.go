package admin

import (
	"skytakeout/global"
	"skytakeout/internal/api/controller"
	"skytakeout/internal/dao"
	"skytakeout/internal/service"
	"skytakeout/middlewares"

	"github.com/gin-gonic/gin"
)

type EmployeeRouter struct {
	service service.IEmployeeService
}

func (er *EmployeeRouter) InitApiRouter(router *gin.RouterGroup) {
	// 拆分鉴权路由
	publicRouter := router.Group("employee")
	privateRouter := router.Group("employee")
	// 私有路由使用jwt验证管理员
	privateRouter.Use(middlewares.VerifyJWTAdminV1())
	// 依赖注入
	er.service = service.NewEmployeeService(
		dao.NewEmployeeDao(global.DB),
	)
	employeeCtl := controller.NewEmployeeController(er.service)
	{
		publicRouter.POST("/login", employeeCtl.Login)
		// 以下都需要jwt登录验证
		privateRouter.POST("/logout", employeeCtl.Logout)
		privateRouter.GET("/session/status", employeeCtl.CheckSessionStatus)
		privateRouter.POST("", employeeCtl.AddEmployee)
		privateRouter.POST("/status/:status", employeeCtl.OnOrOff)
		privateRouter.PUT("/editPassword", employeeCtl.EditPassword)
		privateRouter.PUT("", employeeCtl.UpdateEmployee)
		privateRouter.GET("/page", employeeCtl.PageQuery)
		privateRouter.GET("/:id", employeeCtl.GetById)
	}
}
