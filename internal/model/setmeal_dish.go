package model

// type SetMealDish struct {
// 	Id        uint64  `json:"id"`        // 中间表id
// 	SetmealId uint64  `json:"setmealId"` // 套餐id
// 	DishId    uint64  `json:"dishId"`    // 菜品id
// 	Name      string  `json:"name"`      // 菜品名称冗余字段
// 	Price     float64 `json:"price"`     // 菜品单价冗余字段
// 	Copies    int     `json:"copies"`    // 菜品数量冗余字段
// }

type SetMealDish struct {
	ID        uint64   `gorm:"primaryKey;autoIncrement" json:"id"`
	SetmealID *int64   `json:"setmealId"`
	DishID    *int64   `json:"dishId"`
	Name      string   `gorm:"size:32" json:"name"`
	Price     *float64 `gorm:"type:decimal(10,2)" json:"price"`
	Copies    *int     `json:"copies"`
}

func (s *SetMealDish) TableName() string {
	return "setmeal_dish"
}
