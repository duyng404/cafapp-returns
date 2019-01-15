package gorm

import (
	"cafapp-returns/logger"
)

// OrderStatus statuses of orders
type OrderStatus struct {
	Name        string
	Description string // for frontend
	StatusCode  int    `gorm:"primary_key,unique"`
	ColorCode   string
}

var (
	statusList = []OrderStatus{
		OrderStatus{
			Name:        "incomplete",
			Description: "incomplete",
			StatusCode:  0,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "need-drink",
			Description: "incomplete",
			StatusCode:  1,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "need-destination",
			Description: "incomplete",
			StatusCode:  2,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "need-gus-id",
			Description: "incomplete",
			StatusCode:  3,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "completed",
			Description: "completed but not placed",
			StatusCode:  10,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "placed",
			Description: "placed",
			StatusCode:  20,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "queued",
			Description: "queued",
			StatusCode:  30,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "prepping",
			Description: "prepping",
			StatusCode:  40,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "shipping",
			Description: "shipping",
			StatusCode:  50,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "settled",
			Description: "settled",
			StatusCode:  60,
			ColorCode:   "#ffffff",
		},
		OrderStatus{
			Name:        "general-failure",
			Description: "failed",
			StatusCode:  90,
			ColorCode:   "#ff0000",
		},
		OrderStatus{
			Name:        "unable-to-charge",
			Description: "failed",
			StatusCode:  90,
			ColorCode:   "#ff0000",
		},
		OrderStatus{
			Name:        "abandoned",
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
