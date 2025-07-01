package service

import (
	"context"
	"fmt"
	"skytakeout/common"
	"skytakeout/common/e"
	"skytakeout/common/enum"
	"skytakeout/common/retcode"
	"skytakeout/common/utils"
	"skytakeout/global"
	"skytakeout/internal/api/request"
	"skytakeout/internal/api/response"
	"skytakeout/internal/cache"
	"skytakeout/internal/dao"
	"skytakeout/internal/model"
	"skytakeout/logger"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type IEmployeeService interface {
	Login(context.Context, request.EmployeeLogin) (*response.EmployeeLogin, error)
	Logout(context.Context, string, string) error
	EditPassword(context.Context, request.EmployeeEditPassword) error
	CreateEmployee(ctx context.Context, employee request.EmployeeDTO) error
	PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error)
	SetStatus(ctx context.Context, id uint64, status int) error
	UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error
	GetById(ctx context.Context, id uint64) (*model.Employee, error)
}
type EmployeeImpl struct {
	repo *dao.EmployeeDao
}

func NewEmployeeService(repo *dao.EmployeeDao) IEmployeeService {
	return &EmployeeImpl{repo: repo}
}

// 新增员工
func (ei *EmployeeImpl) CreateEmployee(ctx context.Context, employeeDTO request.EmployeeDTO) error {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "CreateEmployee")
	defer span.End()
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
	fmt.Println(entity.Password)
	if err != nil {
		logger.Logger(ctx2).Error("utils.SetPassword failed", zap.Error(err))
		return err
	}
	// 新增用户
	err = ei.repo.Insert(ctx2, entity)
	if err != nil {
		logger.Logger(ctx2).Error("repo.Insert failed", zap.Error(err))
		return err
	}
	return nil
}

// 登录业务
func (ei *EmployeeImpl) Login(ctx context.Context, employeeLogin request.EmployeeLogin) (*response.EmployeeLogin, error) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "Login")
	defer span.End()
	// 1.查询用户是否存在
	employee, err := ei.repo.GetByUserName(ctx2, employeeLogin.UserName)
	if err != nil || employee == nil {
		logger.Logger(ctx2).Error("repo.GetByUserName failed", zap.Error(err))
		return nil, retcode.NewError(e.ErrorAccountNotFound, e.GetMsg(e.ErrorAccountNotFound))
	}
	// 2.校验密码
	// password := utils.MD5V(employeeLogin.Password, "", 0)
	err = utils.CheckPassword(employee.Password, employeeLogin.Password)
	// if password != employee.Password {
	// 	return nil, retcode.NewError(e.ErrorPasswordError, e.GetMsg(e.ErrorPasswordError))
	// }
	if err != nil {
		logger.Logger(ctx2).Error("utils.CheckPassword failed", zap.Error(err))
		return nil, retcode.NewError(e.ErrorPasswordError, e.GetMsg(e.ErrorPasswordError))
	}
	// 3.校验状态
	if employee.Status == enum.DISABLE {
		logger.Logger(ctx2).Error("Status.DISABLE failed", zap.Error(err))
		return nil, retcode.NewError(e.ErrorAccountLOCKED, e.GetMsg(e.ErrorAccountLOCKED))
	}
	// 4. 生成token
	jwtConfig := global.Config.Jwt.Admin
	accessToken, refreshToken, err := utils.GenerateTokenV1(employee.Id, employeeLogin.UserName, jwtConfig.Secret)
	if err != nil {
		logger.Logger(ctx2).Error("utils.GenerateToken failed", zap.Error(err))
		return nil, err
	}
	// 5. token存入redis
	err = cache.StoreUserAToken(ctx2, accessToken, employeeLogin.UserName)
	if err != nil {
		logger.Logger(ctx2).Error("cache.StoreUserAToken failed", zap.Error(err))
		return nil, err
	}
	err = cache.StoreUserRToken(ctx2, refreshToken, employeeLogin.UserName)
	if err != nil {
		logger.Logger(ctx2).Error("cache.StoreUserRToken failed", zap.Error(err))
		return nil, err
	}
	// 6.构造返回数据
	resp := response.EmployeeLogin{
		Id:           employee.Id,
		Name:         employee.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserName:     employee.Username,
	}
	return &resp, nil
}

func (ei *EmployeeImpl) Logout(ctx context.Context, userName, accessToken string) error {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "Logout")
	defer span.End()
	// 执行退出操作
	if accessToken != "" {
		err := cache.DeleteUserAToken(ctx2, userName)
		if err != nil {
			logger.Logger(ctx2).Error("cache.DeleteUserAToken failed", zap.Error(err))
			return err
		}
	}
	logger.Logger(ctx2).Info("token已经清除,登录状态失效")
	return nil
}

// 设置用户状态
func (ei *EmployeeImpl) SetStatus(ctx context.Context, id uint64, status int) error {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "SetStatus")
	defer span.End()
	entity := model.Employee{Id: id, Status: status}
	err := ei.repo.UpdateStatus(ctx2, entity)
	if err != nil {
		logger.Logger(ctx2).Error("repo.UpdateStatus failed", zap.Error(err))
		return err
	}
	return nil
}

// 修改密码
func (ei *EmployeeImpl) EditPassword(ctx context.Context, employeeEdit request.EmployeeEditPassword) error {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "SetStatus")
	defer span.End()
	// 1.获取员工信息
	employee, err := ei.repo.GetById(ctx2, employeeEdit.EmpId)
	if err != nil {
		logger.Logger(ctx2).Error("repo.GetById failed", zap.Error(err))
		return err
	}
	// 校验用户老密码
	if employee == nil {
		logger.Logger(ctx2).Error("repo.GetById failed", zap.Error(err))
		return retcode.NewError(e.ErrorAccountNotFound, e.GetMsg(e.ErrorAccountNotFound))
	}
	// oldHashPassword := utils.MD5V(employeeEdit.OldPassword, "", 0)
	// if employee.Password != oldHashPassword {
	// 	return retcode.NewError(e.ErrorPasswordError, e.GetMsg(e.ErrorPasswordError))
	// }
	err = utils.CheckPassword(employee.Password, employeeEdit.OldPassword)
	if err != nil {
		logger.Logger(ctx2).Error("utils.CheckPassword failed", zap.Error(err))
		return retcode.NewError(e.ErrorPasswordError, e.GetMsg(e.ErrorPasswordError))
	}
	// 修改员工密码
	// newHashPassword := utils.MD5V(employeeEdit.NewPassword, "", 0) // 使用新密码生成哈希值
	newHashPassword, err := utils.SetPassword(employeeEdit.NewPassword)
	if err != nil {
		logger.Logger(ctx2).Error("utils.SetPassword failed", zap.Error(err))
		return err
	}
	err = ei.repo.Update(ctx2, model.Employee{
		Id:       employeeEdit.EmpId,
		Password: newHashPassword,
	})
	if err != nil {
		logger.Logger(ctx2).Error("repo.Update failed", zap.Error(err))
		return err
	}
	return nil
}

// 更新员工业务
func (ei *EmployeeImpl) UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "UpdateEmployee")
	defer span.End()
	// 构建model实体进行更新
	err := ei.repo.Update(ctx2, model.Employee{
		Id:       dto.Id,
		Username: dto.UserName,
		Name:     dto.Name,
		Phone:    dto.Phone,
		Sex:      dto.Sex,
		IdNumber: dto.IdNumber,
	})
	if err != nil {
		logger.Logger(ctx2).Error("repo.Update failed", zap.Error(err))
		return err
	}
	return nil
}

// 员工分页查询业务
func (ei *EmployeeImpl) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "PageQuery")
	defer span.End()
	// 分页查询
	pageResult, err := ei.repo.PageQuery(ctx2, dto)
	if err != nil {
		logger.Logger(ctx2).Error("repo.PageQuery failed", zap.Error(err))
		return nil, err
	}
	// 屏蔽敏感信息
	if employees, ok := pageResult.Records.([]model.Employee); ok {
		// 替换敏感信息
		for key, _ := range employees {
			employees[key].Password = "****"
			employees[key].IdNumber = "****"
			employees[key].Phone = "****"
		}
		// 重新赋值
		pageResult.Records = employees
	}
	return pageResult, nil
}

// 根据id获取员工id
func (ei *EmployeeImpl) GetById(ctx context.Context, id uint64) (*model.Employee, error) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "GetById")
	defer span.End()
	employee, err := ei.repo.GetById(ctx2, id)
	if err != nil {
		logger.Logger(ctx2).Error("repo.GetById failed", zap.Error(err))
		return nil, err
	}
	employee.Password = "***"
	return employee, err
}
