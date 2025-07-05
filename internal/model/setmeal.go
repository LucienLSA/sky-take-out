package model

import (
	"skytakeout/common/enum"
	"time"

	"gorm.io/gorm"
)

// type SetMeal struct {
// 	Id          uint64    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"` // 主键id
// 	CategoryId  uint64    `json:"category_id"`                         // 分类id
// 	Name        string    `json:"name"`                                // 套餐名称
// 	Price       float64   `json:"price"`                               // 套餐单价
// 	Status      int       `json:"status"`                              // 套餐状态
// 	Description string    `json:"description"`                         // 套餐描述
// 	Image       string    `json:"image"`                               // 套餐图片
// 	CreateTime  time.Time `json:"create_time"`                         // 创建时间
// 	UpdateTime  time.Time `json:"update_time"`                         // 更新时间
// 	CreateUser  uint64    `json:"create_user"`                         // 创建用户
// 	UpdateUser  uint64    `json:"update_user"`                         // 更新用户
// }

type SetMeal struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID  int64      `gorm:"not null" json:"categoryId"`
	Name        string     `gorm:"size:32;not null;uniqueIndex" json:"name"`
	Price       float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Status      *int       `gorm:"default:1" json:"status"` // 0:停售 1:起售
	Description string     `gorm:"size:255" json:"description"`
	Image       string     `gorm:"size:255" json:"image"`
	CreateTime  *time.Time `json:"createTime"`
	UpdateTime  *time.Time `json:"updateTime"`
	CreateUser  *int64     `json:"createUser"`
	UpdateUser  *int64     `json:"updateUser"`
}

func (e *SetMeal) BeforeCreate(tx *gorm.DB) error {
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

func (e *SetMeal) BeforeUpdate(tx *gorm.DB) error {
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

func (e *SetMeal) TableName() string {
	return "setmeal"
}
