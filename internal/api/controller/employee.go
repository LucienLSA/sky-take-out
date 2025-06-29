package controller

import (
	"skytakeout/common/retcode"
	"skytakeout/global"
	"skytakeout/internal/api/request"
	"skytakeout/internal/service"

	"github.com/gin-gonic/gin"
)

// 控制器依赖于服务层的抽象，而不是具体实现
type EmployeeController struct {
	service service.IEmployeeService
}

// 构造函数
// 参数 employeeService 是 IEmployeeService 接口类型
// 函数返回一个 *EmployeeController 指针
// 在创建控制器时，将传入的服务实例赋值给控制器的 service 字段
func NewEmployeeController(employeeService service.IEmployeeService) *EmployeeController {
	return &EmployeeController{service: employeeService}
}

// AddEmployee 新增员工
func (ec *EmployeeController) AddEmployee(ctx *gin.Context) {
	var request request.EmployeeDTO
	var err error
	err = ctx.Bind(&request)
	if err != nil {
		global.Log.Error(ctx, "AddEmployee Error: err=%s", err.Error())
		retcode.Fatal(ctx, err, "")
		return
	}
	err = ec.service.CreateEmployee(ctx, request)
	if err != nil {
		global.Log.Error(ctx, "AddEmployee  Error: err=%s", err.Error())
		retcode.Fatal(ctx, err, "")
		return
	}
	// 正确输出
	retcode.OK(ctx, "")
}

// Login 员工登录
func (ec *EmployeeController) Login(ctx *gin.Context) {
	employeeLogin := request.EmployeeLogin{}
	err := ctx.Bind(&employeeLogin)
	if err != nil {
		global.Log.Error(ctx, "EmployeeController login 解析失败")
		retcode.Fatal(ctx, err, "")
		return
	}
	resp, err := ec.service.Login(ctx, employeeLogin)
	if err != nil {
		global.Log.Error(ctx, "EmployeeController login Error: err=%s", err.Error())
		retcode.Fatal(ctx, err, "")
		return
	}
	retcode.OK(ctx, resp)
}

// Logout 员工退出
func (ec *EmployeeController) Logout(ctx *gin.Context) {
	var err error
	err = ec.service.Logout(ctx)
	if err != nil {
		global.Log.Error(ctx, "EmployeeController login Error: err=%s", err.Error())
		retcode.Fatal(ctx, err, "")
		return
	}
	retcode.OK(ctx, "")
}
