package dao

import (
	"context"
	"skytakeout/common/e"
	"skytakeout/common/retcode"
	"skytakeout/global"
	"skytakeout/internal/model"

	"gorm.io/gorm"
)

type EmployeeDao struct {
	db *gorm.DB
}

func NewEmployeeDao(db *gorm.DB) *EmployeeDao {
	return &EmployeeDao{db: db}
}
func (d *EmployeeDao) GetByUserName(ctx context.Context, userName string) (*model.Employee, error) {
	var employee model.Employee
	err := d.db.WithContext(ctx).Where("username=?", userName).First(&employee).Error
	if err != nil {
		global.Log.Error(ctx, "EmployeeDao.GetByUserName failed, err: %v", err)
		return nil, retcode.NewError(e.MysqlERR, "Get employee failed")
	}
	return &employee, nil
}
func (d *EmployeeDao) Insert(ctx context.Context, entity model.Employee) error {
	err := d.db.WithContext(ctx).Create(&entity).Error
	if err != nil {
		global.Log.Error(ctx, "EmployeeDao.Insert failed, err: %v", err)
		return retcode.NewError(e.MysqlERR, "Create employee failed")
	}
	return nil
}
