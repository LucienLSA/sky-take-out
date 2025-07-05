package model

import (
	"skytakeout/common/enum"
	"time"

	"gorm.io/gorm"
)

// type Category struct {
// 	Id         uint64    `json:"id"`
// 	Type       int       `json:"type"`
// 	Name       string    `json:"name"`
// 	Sort       int       `json:"sort"`
// 	Status     int       `json:"status"`
// 	CreateTime time.Time `json:"createTime"`
// 	UpdateTime time.Time `json:"updateTime"`
// 	CreateUser uint64    `json:"createUser"`
// 	UpdateUser uint64    `json:"updateUser"`
// }

type Category struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Type       *int       `json:"type"` // 1 菜品分类 2 套餐分类
	Name       string     `gorm:"size:32;not null;uniqueIndex" json:"name"`
	Sort       int        `gorm:"default:0;not null" json:"sort"`
	Status     *int       `json:"status"` // 0:禁用，1:启用
	CreateTime *time.Time `json:"createTime"`
	UpdateTime *time.Time `json:"updateTime"`
	CreateUser *int64     `json:"createUser"`
	UpdateUser *int64     `json:"updateUser"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	// 自动填充 创建时间、创建人、更新时间、更新用户
	now := time.Now()
	c.CreateTime = &now
	c.UpdateTime = &now
	// 从上下文获取用户信息
	value := tx.Statement.Context.Value(enum.CurrentId)
	if uid, ok := value.(int64); ok {
		c.CreateUser = &uid
		c.UpdateUser = &uid
	}
	return nil
}

func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	// 在更新记录千自动填充更新时间
	now := time.Now()
	c.UpdateTime = &now
	// 从上下文获取用户信息
	value := tx.Statement.Context.Value(enum.CurrentId)
	if uid, ok := value.(int64); ok {
		c.UpdateUser = &uid
	}
	return nil
}
