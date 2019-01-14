package gorm

import (
	"cafapp-returns/logger"

	"github.com/jinzhu/gorm"
)

// OrderStatus statuses of orders
type OrderStatus struct {
	gorm.Model
	Name        string
	Description string // for frontend
	StatusCode  int
	ColorCode   string
}

var (
	statusList = []OrderStatus{
		OrderStatus{
			Name:        "Incomplete",
			Description: "incomplete",
			StatusCode:  0,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "Completed",
			Description: "completed but not placed",
			StatusCode:  10,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "Placed",
			Description: "placed",
			StatusCode:  20,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "Queued",
			Description: "queued",
			StatusCode:  30,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "Prepping",
			Description: "prepping",
			StatusCode:  40,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "Shipping",
			Description: "shipping",
			StatusCode:  50,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "Settled",
			Description: "settled",
			StatusCode:  60,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "General Failure",
			Description: "failed",
			StatusCode:  90,
			ColorCode:   "#ff0000",
		},
		OrderStatus{
			Name:        "Abandoned",
			Description: "failed",
			StatusCode:  91,
			ColorCode:   "#ff0000",
		},
	}
)

// PopulateByCode query db by code
func (s *OrderStatus) PopulateByCode(code int) error {
	return DB.Where("status_code = ?", code).Last(&s).Error
}

// FirstOrCreate create if not exist
func (s *OrderStatus) FirstOrCreate() error {
	var tmp OrderStatus
	if err := tmp.PopulateByCode(s.StatusCode); err != nil {
		return DB.Create(&s).Error
	}
	return nil
}

// create all statuses
func initOrderStatuses() error {
	for i := range statusList {
		err := statusList[i].FirstOrCreate()
		if err != nil {
			logger.Error(err)
		}
	}
	return nil
}
