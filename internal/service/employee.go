package service

import (
	"context"
	"skytakeout/common/e"
	"skytakeout/common/enum"
	"skytakeout/common/retcode"
	"skytakeout/common/utils"
	"skytakeout/global"
	"skytakeout/internal/api/request"
	"skytakeout/internal/api/response"
	"skytakeout/internal/dao"
	"skytakeout/internal/model"
)

type IEmployeeService interface {
	Login(context.Context, request.EmployeeLogin) (*response.EmployeeLogin, error)
	Logout(ctx context.Context) error
	// EditPassword(context.Context, request.EmployeeEditPassword) error
	CreateEmployee(ctx context.Context, employee request.EmployeeDTO) error
	// PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error)
	// SetStatus(ctx context.Context, id uint64, status int) error
	// UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error
	// GetById(ctx context.Context, id uint64) (*model.Employee, error)
}
type EmployeeImpl struct {
	repo *dao.EmployeeDao
}

func NewEmployeeService(repo *dao.EmployeeDao) IEmployeeService {
	return &EmployeeImpl{repo: repo}
}

// 新增员工
func (ei *EmployeeImpl) CreateEmployee(ctx context.Context, employeeDTO request.EmployeeDTO) error {
	var err error
	// 1.新增员工,构建员工基础信息
	entity := model.Employee{
		Id:       employeeDTO.Id,
		IdNumber: employeeDTO.IdNumber,
		Name:     employeeDTO.Name,
		Phone:    employeeDTO.Phone,
		Sex:      employeeDTO.Sex,
		Username: employeeDTO.UserName,
	}
	// 新增用户为启用状态
	entity.Status = enum.ENABLE
	// 新增用户初始密码为123456
	// entity.Password = utils.MD5V("123456", "", 0)
	entity.Password, err = utils.SetPassword("123456")
	if err != nil {
		global.Log.Error(ctx, "utils.SetPassword failed, err: %v", err)
		return err
	}
	// 新增用户
	err = ei.repo.Insert(ctx, entity)
	if err != nil {
		global.Log.Error(ctx, "EmployeeImpl.CreateEmployee failed, err: %v", err)
		return err
	}
	return nil
}

// 登录业务
func (ei *EmployeeImpl) Login(ctx context.Context, employeeLogin request.EmployeeLogin) (*response.EmployeeLogin, error) {
	// 1.查询用户是否存在
	employee, err := ei.repo.GetByUserName(ctx, employeeLogin.UserName)
	if err != nil || employee == nil {
		return nil, retcode.NewError(e.ErrorAccountNotFound, e.GetMsg(e.ErrorAccountNotFound))
	}
	// 2.校验密码
	password := utils.MD5V(employeeLogin.Password, "", 0)
	if password != employee.Password {
		return nil, retcode.NewError(e.ErrorPasswordError, e.GetMsg(e.ErrorPasswordError))
	}
	// 3.校验状态
	if employee.Status == enum.DISABLE {
		return nil, retcode.NewError(e.ErrorAccountLOCKED, e.GetMsg(e.ErrorAccountLOCKED))
	}
	// 生成Token
	jwtConfig := global.Config.Jwt.Admin
	token, err := utils.GenerateToken(employee.Id, jwtConfig.Name, jwtConfig.Secret)
	if err != nil {
		global.Log.Error(ctx, "EmployeeImpl.Login failed, err: %v", err)
		return nil, err
	}
	// 4.构造返回数据
	resp := response.EmployeeLogin{
		Id:       employee.Id,
		Name:     employee.Name,
		Token:    token,
		UserName: employee.Username,
	}
	return &resp, nil
}

func (ei *EmployeeImpl) Logout(ctx context.Context) error {
	// TODO 后续扩展为单点登录模式。 1.获取上下文中当前用户
	// 2.如果是单点登录的话执行推出操作
	return nil
}
