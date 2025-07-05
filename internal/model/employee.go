package model

import (
	"skytakeout/common/enum"
	"time"

	"gorm.io/gorm"
)

//	type Employee struct {
//		Id         uint64    `json:"id"`
//		Username   string    `json:"username"`
//		Name       string    `json:"name"`
//		Password   string    `json:"password"`
//		Phone      string    `json:"phone"`
//		Sex        string    `json:"sex"`
//		IdNumber   string    `json:"idNumber"`
//		Status     int       `json:"status"`
//		CreateTime time.Time `json:"createTime"`
//		UpdateTime time.Time `json:"updateTime"`
//		CreateUser uint64    `json:"createUser"`
//		UpdateUser uint64    `json:"updateUser"`
//	}
type Employee struct {
	Id         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string     `gorm:"size:32;not null" json:"name"`
	Username   string     `gorm:"size:32;not null;uniqueIndex" json:"username"`
	Password   string     `gorm:"size:64;not null" json:"password"`
	Phone      string     `gorm:"size:11;not null" json:"phone"`
	Sex        string     `gorm:"size:2;not null" json:"sex"`
	IdNumber   string     `gorm:"size:18;not null" json:"idNumber"`
	Status     int        `gorm:"default:1;not null" json:"status"` // 0:禁用，1:启用
	CreateTime *time.Time `json:"createTime"`
	UpdateTime *time.Time `json:"updateTime"`
	CreateUser *int64     `json:"createUser"`
	UpdateUser *int64     `json:"updateUser"`
}

func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	// 自动填充 创建时间、创建人、更新时间、更新用户
	now := time.Now()
	e.CreateTime = &now
	e.UpdateTime = &now
	// 从上下文获取用户信息
	value := tx.Statement.Context.Value(enum.CurrentId)
	if uid, ok := value.(int64); ok {
		e.CreateUser = &uid
		e.UpdateUser = &uid
	}
	return nil
}

func (e *Employee) BeforeUpdate(tx *gorm.DB) error {
	// 在更新记录千自动填充更新时间
	now := time.Now()
	e.UpdateTime = &now
	// 从上下文获取用户信息
	value := tx.Statement.Context.Value(enum.CurrentId)
	if uid, ok := value.(int64); ok {
		e.UpdateUser = &uid
	}
	return nil
}

func (e *Employee) AfterFind(tx *gorm.DB) error {
	// 格式化当前日期
	// e.CreateTime.Format(time.DateOnly)
	// e.CreateTime.Format(time.DateTime)
	return nil
}
