package model

// type DishFlavor struct {
// 	Id     uint64 `json:"id"`      //口味id
// 	DishId uint64 `json:"dish_id"` //菜品id
// 	Name   string `json:"name"`    //口味主题 温度|甜度|辣度
// 	Value  string `json:"value"`   //口味信息 可多个
// }

type DishFlavor struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	DishID int64  `gorm:"not null" json:"dishId"`
	Name   string `gorm:"size:32" json:"name"`
	Value  string `gorm:"size:255" json:"value"`
}

func (d *DishFlavor) TableName() string {
	return "dish_flavor"
}
