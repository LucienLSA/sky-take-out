package model

import "time"

type Orders struct {
	ID                    int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Number                string     `gorm:"size:50" json:"number"`
	Status                int        `gorm:"default:1;not null" json:"status"` // 1待付款 2待接单 3已接单 4派送中 5已完成 6已取消 7退款
	UserID                int64      `gorm:"not null" json:"userId"`
	AddressBookID         int64      `gorm:"not null" json:"addressBookId"`
	OrderTime             time.Time  `gorm:"not null" json:"orderTime"`
	CheckoutTime          *time.Time `json:"checkoutTime"`
	PayMethod             int        `gorm:"default:1;not null" json:"payMethod"` // 1微信,2支付宝
	PayStatus             int8       `gorm:"default:0;not null" json:"payStatus"` // 0未支付 1已支付 2退款
	Amount                float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Remark                string     `gorm:"size:100" json:"remark"`
	Phone                 string     `gorm:"size:11" json:"phone"`
	Address               string     `gorm:"size:255" json:"address"`
	UserName              string     `gorm:"size:32" json:"userName"`
	Consignee             string     `gorm:"size:32" json:"consignee"`
	CancelReason          string     `gorm:"size:255" json:"cancelReason"`
	RejectionReason       string     `gorm:"size:255" json:"rejectionReason"`
	CancelTime            *time.Time `json:"cancelTime"`
	EstimatedDeliveryTime *time.Time `json:"estimatedDeliveryTime"`
	DeliveryStatus        int8       `gorm:"default:1;not null" json:"deliveryStatus"` // 1立即送出 0选择具体时间
	DeliveryTime          *time.Time `json:"deliveryTime"`
	PackAmount            *int       `json:"packAmount"`
	TablewareNumber       *int       `json:"tablewareNumber"`
	TablewareStatus       int8       `gorm:"default:1;not null" json:"tablewareStatus"` // 1按餐量提供 0选择具体数量
}
