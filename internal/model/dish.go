package model

import (
	"skytakeout/common/enum"
	"time"

	"gorm.io/gorm"
)

//	type Dish struct {
//		Id          uint64    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
//		Name        string    `json:"name"`
//		CategoryId  uint64    `json:"categoryId"`
//		Price       float64   `json:"price"`
//		Image       string    `json:"image"`
//		Description string    `json:"description"`
//		Status      int       `json:"status"`
//		CreateTime  time.Time `json:"createTime"`
//		UpdateTime  time.Time `json:"updateTime"`
//		CreateUser  uint64    `json:"createUser"`
//		UpdateUser  uint64    `json:"updateUser"`
//		// 一对多
//		Flavors []DishFlavor `json:"flavors"`
//	}
type Dish struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"size:32;not null;uniqueIndex" json:"name"`
	CategoryID  int64      `gorm:"not null" json:"categoryId"`
	Price       *float64   `gorm:"type:decimal(10,2)" json:"price"`
	Image       string     `gorm:"size:255" json:"image"`
	Description string     `gorm:"size:255" json:"description"`
	Status      *int       `gorm:"default:1" json:"status"` // 0 停售 1 起售
	CreateTime  *time.Time `json:"createTime"`
	UpdateTime  *time.Time `json:"updateTime"`
	CreateUser  *int64     `json:"createUser"`
	UpdateUser  *int64     `json:"updateUser"`
}

func (d *Dish) BeforeCreate(tx *gorm.DB) error {
	// 自动填充 创建时间、创建人、更新时间、更新用户
	now := time.Now()
	d.CreateTime = &now
	d.UpdateTime = &now
	// 从上下文获取用户信息
	value := tx.Statement.Context.Value(enum.CurrentId)
	if uid, ok := value.(int64); ok {
		d.CreateUser = &uid
		d.UpdateUser = &uid
	}
	return nil
}

func (d *Dish) BeforeUpdate(tx *gorm.DB) error {
	// 在更新记录千自动填充更新时间
	now := time.Now()
	d.UpdateTime = &now
	// 从上下文获取用户信息
	value := tx.Statement.Context.Value(enum.CurrentId)
	if uid, ok := value.(int64); ok {
		d.UpdateUser = &uid
	}
	return nil
}

func (e *Dish) TableName() string {
	return "dish"
}
