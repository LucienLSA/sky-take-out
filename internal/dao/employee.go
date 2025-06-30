package dao

import (
	"context"
	"skytakeout/common"
	"skytakeout/common/e"
	"skytakeout/common/retcode"
	"skytakeout/global"
	"skytakeout/internal/api/request"
	"skytakeout/internal/model"

	"skytakeout/logger"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EmployeeDao struct {
	db *gorm.DB
}

func NewEmployeeDao(db *gorm.DB) *EmployeeDao {
	return &EmployeeDao{
		db: db,
	}
}

// 根据员工名获取员工
func (d *EmployeeDao) GetByUserName(ctx context.Context, userName string) (*model.Employee, error) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "GetByUserName")
	defer span.End()
	var employee model.Employee
	err := d.db.WithContext(ctx2).Where("username=?", userName).First(&employee).Error
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeDao.GetByUserName failed", zap.Error(err))
		return nil, retcode.NewError(e.MysqlERR, "Get employee failed")
	}
	return &employee, nil
}

// 新增员工
func (d *EmployeeDao) Insert(ctx context.Context, entity model.Employee) error {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "Insert")
	defer span.End()
	err := d.db.WithContext(ctx2).Create(&entity).Error
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeDao.Insert failed", zap.Error(err))
		return retcode.NewError(e.MysqlERR, "Create employee failed")
	}
	return nil
}

// UpdateStatus 动态更新包括零值
func (d *EmployeeDao) UpdateStatus(ctx context.Context, employee model.Employee) error {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "UpdateStatus")
	defer span.End()
	err := d.db.WithContext(ctx2).Model(&model.Employee{}).Where("id = ?",
		employee.Id).Update("status", employee.Status).Error
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeDao.UpdateStatus failed", zap.Error(err))
		return retcode.NewError(e.MysqlERR, "update employee failed")
	}
	return nil
}

// 根据员工id获取员工
func (d *EmployeeDao) GetById(ctx context.Context, id uint64) (*model.Employee, error) {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "GetById")
	defer span.End()
	var employee model.Employee
	err := d.db.WithContext(ctx2).Where("id=?", id).First(&employee).Error
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeDao.GetById failed", zap.Error(err))
		return nil, retcode.NewError(e.MysqlERR, "Get employee failed")
	}
	return &employee, nil
}

// 更新员工信息
func (d *EmployeeDao) Update(ctx context.Context, employee model.Employee) error {
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "Update")
	defer span.End()
	err := d.db.WithContext(ctx2).Model(&employee).Updates(employee).Error
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeDao.Update failed", zap.Error(err))
		return retcode.NewError(e.MysqlERR, "Update employee failed")
	}
	return nil
}

func (d *EmployeeDao) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error) {
	// 分页查询 select count(*) from employee where name = ? limit x,y
	tracer := otel.Tracer(global.ServiceName)
	ctx2, span := tracer.Start(ctx, "PageQuery")
	defer span.End()
	var result common.PageResult
	var employeeList []model.Employee
	var err error
	// 动态拼接
	query := d.db.WithContext(ctx2).Model(&model.Employee{})
	if dto.Name != "" {
		query = query.Where("name LIKE ?", "%"+dto.Name+"%")
	}
	// 计算总数
	if err = query.Count(&result.Total).Error; err != nil {
		logger.Logger(ctx2).Error("EmployeeDao.PageQuery Count failed", zap.Error(err))
		return nil, retcode.NewError(e.MysqlERR, "Get employee List failed")
	}
	// 分页查询
	err = query.Scopes(result.Paginate(&dto.Page, &dto.PageSize)).Find(&employeeList).Error
	if err != nil {
		logger.Logger(ctx2).Error("EmployeeDao.PageQuery List failed", zap.Error(err))
		return nil, retcode.NewError(e.MysqlERR, "Get employee List failed")
	}
	result.Records = employeeList
	return &result, nil
}
