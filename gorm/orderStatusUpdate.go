package gorm

import "github.com/jinzhu/gorm"

// OrderStatusUpdate keeps track of when which order changes status and what status it is
type OrderStatusUpdate struct {
	gorm.Model
	OrderUUID  string
	StatusCode int
}

// CreateOrderStatusUpdate whenever an order changes status, we record it so we know the time
func CreateOrderStatusUpdate(uuid string, code int) error {
	tmp := OrderStatusUpdate{
		OrderUUID:  uuid,
		StatusCode: code,
	}
	return DB.Create(&tmp).Error
}
