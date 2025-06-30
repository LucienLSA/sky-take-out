package dao

import (
	"context"
	"skytakeout/common"
	"skytakeout/common/e"
	"skytakeout/common/retcode"
	"skytakeout/global"
	"skytakeout/internal/api/request"
	"skytakeout/internal/model"

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
	var employee model.Employee
	err := d.db.WithContext(ctx).Where("username=?", userName).First(&employee).Error
	if err != nil {
		global.Log.Error(ctx, "EmployeeDao.GetByUserName failed, err: %v", err)
		return nil, retcode.NewError(e.MysqlERR, "Get employee failed")
	}
	return &employee, nil
}

// 新增员工
func (d *EmployeeDao) Insert(ctx context.Context, entity model.Employee) error {
	err := d.db.WithContext(ctx).Create(&entity).Error
	if err != nil {
		global.Log.Error(ctx, "EmployeeDao.Insert failed, err: %v", err)
		return retcode.NewError(e.MysqlERR, "Create employee failed")
	}
	return nil
}

// UpdateStatus 动态更新包括零值
func (d *EmployeeDao) UpdateStatus(ctx context.Context, employee model.Employee) error {
	err := d.db.WithContext(ctx).Model(&model.Employee{}).Where("id = ?",
		employee.Id).Update("status", employee.Status).Error
	if err != nil {
		global.Log.Error(ctx, "EmployeeDao.UpdateStatus failed, err: %v", err)
		return retcode.NewError(e.MysqlERR, "update employee failed")
	}
	return nil
}

// 根据员工id获取员工
func (d *EmployeeDao) GetById(ctx context.Context, id uint64) (*model.Employee, error) {
	var employee model.Employee
	err := d.db.WithContext(ctx).Where("id=?", id).First(&employee).Error
	if err != nil {
		global.Log.Error(ctx, "EmployeeDao.GetById failed, err: %v", err)
		return nil, retcode.NewError(e.MysqlERR, "Get employee failed")
	}
	return &employee, nil
}

// 更新员工信息
func (d *EmployeeDao) Update(ctx context.Context, employee model.Employee) error {
	err := d.db.WithContext(ctx).Model(&employee).Updates(employee).Error
	if err != nil {
		global.Log.Error(ctx, "EmployeeDao.Update failed, err: %v", err)
		return retcode.NewError(e.MysqlERR, "Update employee failed")
	}
	return nil
}

func (d *EmployeeDao) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error) {
	// 分页查询 select count(*) from employee where name = ? limit x,y
	var result common.PageResult
	var employeeList []model.Employee
	var err error
	// 动态拼接
	query := d.db.WithContext(ctx).Model(&model.Employee{})
	if dto.Name != "" {
		query = query.Where("name LIKE ?", "%"+dto.Name+"%")
	}
	// 计算总数
	if err = query.Count(&result.Total).Error; err != nil {
		global.Log.Error(ctx, "EmployeeDao.PageQuery Count failed, err: %v", err)
		return nil, retcode.NewError(e.MysqlERR, "Get employee List failed")
	}
	// 分页查询
	err = query.Scopes(result.Paginate(&dto.Page, &dto.PageSize)).Find(&employeeList).Error
	if err != nil {
		global.Log.Error(ctx, "EmployeeDao.PageQuery List failed, err: %v", err)
		return nil, retcode.NewError(e.MysqlERR, "Get employee List failed")
	}
	result.Records = employeeList
	return &result, nil
}
