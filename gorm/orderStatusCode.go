package gorm

import (
	"cafapp-returns/logger"

	"github.com/jinzhu/gorm"
)

// OrderStatusCode statuses of orders
type OrderStatusCode struct {
	StatusCode  int `gorm:"primary_key;unique"`
	Name        string
	DisplayName string // for frontend
	ColorCode   string
}

const (
	OrderStatusIncomplete = 1
	OrderStatusNeedInfo   = 2

	OrderStatusFinalized         = 10
	OrderStatusInsufficientFunds = 11

	OrderStatusPlaced = 20

	OrderStatusQueued   = 30 // unused
	OrderStatusRequeued = 31 // unused

	OrderStatusPrepping    = 40
	OrderStatusShipping    = 50
	OrderStatusApproaching = 51
	OrderStatusDelivered   = 60

	OrderStatusGeneralFailure = 90
	OrderStatusDeliveryFailed = 92
	OrderStatusFailedToCharge = 93
)

var (
	statusList = []OrderStatusCode{
		OrderStatusCode{
			Name:        "incomplete",
			DisplayName: "Incomplete",
			StatusCode:  OrderStatusIncomplete,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "need-info",
			DisplayName: "Need Info",
			StatusCode:  OrderStatusNeedInfo,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "finalized",
			DisplayName: "Finalized",
			StatusCode:  OrderStatusFinalized,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "insufficient-funds",
			DisplayName: "Insufficient Funds",
			StatusCode:  OrderStatusInsufficientFunds,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "placed",
			DisplayName: "Placed",
			StatusCode:  OrderStatusPlaced,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "queued",
			DisplayName: "Queued",
			StatusCode:  OrderStatusQueued,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "requeued",
			DisplayName: "Re-queued",
			StatusCode:  OrderStatusRequeued,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "prepping",
			DisplayName: "Prepping",
			StatusCode:  OrderStatusPrepping,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "shipping",
			DisplayName: "Out For Delivery",
			StatusCode:  OrderStatusShipping,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "approaching",
			DisplayName: "Arrived At Door",
			StatusCode:  OrderStatusApproaching,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "delivered",
			DisplayName: "Delivered",
			StatusCode:  OrderStatusDelivered,
			ColorCode:   "#ffffff",
		},
		OrderStatusCode{
			Name:        "general-failure",
			DisplayName: "Failed",
			StatusCode:  OrderStatusGeneralFailure,
			ColorCode:   "#ff0000",
		},
		OrderStatusCode{
			Name:        "delivery-failed",
			DisplayName: "Delivery Failed",
			StatusCode:  OrderStatusDeliveryFailed,
			ColorCode:   "#ff0000",
		},
		OrderStatusCode{
			Name:        "failed-to-charge",
			DisplayName: "Failed To Charge",
			StatusCode:  OrderStatusFailedToCharge,
			ColorCode:   "#ff0000",
		},
	}
)

// PopulateByCode query db by code
func (s *OrderStatusCode) PopulateByCode(code int) error {
	return DB.Where("status_code = ?", code).Last(&s).Error
}

// CreateOrUpdate create if not exist
func (s *OrderStatusCode) CreateOrUpdate() error {
	var tmp OrderStatusCode
	err := tmp.PopulateByCode(s.StatusCode)
	if err == gorm.ErrRecordNotFound {
		return DB.Create(&s).Error
	} else if err != nil {
		return err
	} else {
		return DB.Save(s).Error
	}
}

// create all statuses
func initOrderStatusCodes() error {
	for i := range statusList {
		err := statusList[i].CreateOrUpdate()
		if err != nil {
			logger.Error(err)
		}
	}
	return nil
}
