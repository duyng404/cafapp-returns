package gorm

import "github.com/jinzhu/gorm"

// OrderStatusUpdate keeps track of when which order changes status and what status it is
type OrderStatusUpdate struct {
	gorm.Model
	OrderID    uint
	Order      Order
	StatusCode int `json:"status_code"`
}

// CreateOrderStatusUpdate whenever an order changes status, we record it so we know the time
func CreateOrderStatusUpdate(orderID uint, code int) error {
	tmp := OrderStatusUpdate{
		OrderID:    orderID,
		StatusCode: code,
	}
	return DB.Create(&tmp).Error
}
