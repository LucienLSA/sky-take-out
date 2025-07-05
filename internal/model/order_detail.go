package model

type OrderDetail struct {
	ID         uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string  `gorm:"size:32" json:"name"`
	Image      string  `gorm:"size:255" json:"image"`
	OrderID    int64   `gorm:"not null" json:"orderId"`
	DishID     *int64  `json:"dishId"`
	SetmealID  *int64  `json:"setmealId"`
	DishFlavor string  `gorm:"size:50" json:"dishFlavor"`
	Number     int     `gorm:"default:1;not null" json:"number"`
	Amount     float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
}
