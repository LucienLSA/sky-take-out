package controller

import (
	"skytakeout/common/enum"
	"skytakeout/common/retcode"
	"skytakeout/global"
	"skytakeout/internal/api/request"
	"skytakeout/internal/service"
	"strconv"

	"skytakeout/logger"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
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
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "EmployeeController AddEmployee")
	defer span.End()
	var request request.EmployeeDTO
	var err error
	err = ctx.Bind(&request)
	if err != nil {
		logger.Logger(ctx2).Error("AddEmployee Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	err = ec.service.CreateEmployee(ctx2, request)
	if err != nil {
		logger.Logger(ctx2).Error("AddEmployee Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	// 正确输出
	retcode.OK(ctx, "")
}

// Login 员工登录
func (ec *EmployeeController) Login(ctx *gin.Context) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "EmployeeController Login")
	defer span.End()
	employeeLogin := request.EmployeeLogin{}
	err := ctx.Bind(&employeeLogin)
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeController login Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	resp, err := ec.service.Login(ctx2, employeeLogin)
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeController login Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	retcode.OK(ctx, resp)
}

// Logout 员工退出
func (ec *EmployeeController) Logout(ctx *gin.Context) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "EmployeeController Logout")
	defer span.End()
	employeeLogout := request.EmployeeLogout{}
	err := ctx.Bind(&employeeLogout)
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeController Logout Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	// 获取上下文中当前用户
	userName, _ := ctx.Get(enum.CurrentName)
	accessToken := ctx.GetHeader(global.Config.Jwt.Admin.AccessToken)
	err = ec.service.Logout(ctx2, userName.(string), accessToken)
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeController Logout Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	retcode.OK(ctx, "")
}

// OnOrOff 启用Or禁用员工状态
func (ec *EmployeeController) OnOrOff(ctx *gin.Context) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "EmployeeController OnOrOff")
	defer span.End()
	id, _ := strconv.ParseUint(ctx.Query("id"), 10, 64)
	status, _ := strconv.Atoi(ctx.Param("status"))
	var err error
	err = ec.service.SetStatus(ctx2, id, status)
	if err != nil {
		logger.Logger(ctx2).Error("OnOrOff Status Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	// 更新员工状态
	logger.Logger(ctx2).Info("启用Or禁用员工状态", zap.Uint64("id", id), zap.Int("status", status))
	retcode.OK(ctx, "")
}

// EditPassword 修改密码
func (ec *EmployeeController) EditPassword(ctx *gin.Context) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "EmployeeController EditPassword")
	defer span.End()
	var reqs request.EmployeeEditPassword
	var err error
	err = ctx.Bind(&reqs)
	if err != nil {
		logger.Logger(ctx2).Error("EditPassword Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	// 从上下文获取员工id
	if id, ok := ctx.Get(enum.CurrentId); ok {
		reqs.EmpId = id.(uint64)
	}
	err = ec.service.EditPassword(ctx2, reqs)
	if err != nil {
		logger.Logger(ctx2).Error("EditPassword Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
	}
	retcode.OK(ctx, "")
}

// UpdateEmployee 编辑员工信息
func (ec *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "EmployeeController UpdateEmployee")
	defer span.End()
	var employeeDTO request.EmployeeDTO
	err := ctx.Bind(&employeeDTO)
	if err != nil {
		logger.Logger(ctx2).Error("UpdateEmployee Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	// 修改员工信息
	err = ec.service.UpdateEmployee(ctx2, employeeDTO)
	if err != nil {
		logger.Logger(ctx2).Error("UpdateEmployee Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	retcode.OK(ctx, "")
}

// PageQuery 员工分页查询
func (ec *EmployeeController) PageQuery(ctx *gin.Context) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "EmployeeController PageQuery")
	defer span.End()
	var employeePageQueryDTO request.EmployeePageQueryDTO
	err := ctx.Bind(&employeePageQueryDTO)
	if err != nil {
		logger.Logger(ctx2).Error("AddEmployee invalid params", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	// 进行分页查询
	pageResult, err := ec.service.PageQuery(ctx2, employeePageQueryDTO)
	if err != nil {
		logger.Logger(ctx2).Error("AddEmployee Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	retcode.OK(ctx, pageResult)
}

// GetById 获取员工信息根据id
func (ec *EmployeeController) GetById(ctx *gin.Context) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "EmployeeController GetById")
	defer span.End()
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	employee, err := ec.service.GetById(ctx2, id)
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeCtrl GetById Error", zap.Error(err))
		retcode.Fatal(ctx, err, "")
		return
	}
	retcode.OK(ctx, employee)
}
