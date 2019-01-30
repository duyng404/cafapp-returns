package socket

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
)

// commit queue means move from queue to prep
func handleCommitQueue(commited []int) []*gorm.Order {
	if len(commited) > 0 {
		logger.Info("processing queue -> prep:", commited)
		processed := []*gorm.Order{}
		for _, v := range commited {
			// get order from db
			var o gorm.Order
			err := o.PopulateByID(uint(v))
			if err != nil {
				logger.Error("error getting order from db")
				return []*gorm.Order{}
			}
			// set a new status
			err = o.SetStatusTo(gorm.OrderStatusPrepping)
			if err != nil {
				logger.Error("error setting new status")
				return []*gorm.Order{}
			}
			// accumulate to processed
			processed = append(processed, &o)
		}
		return processed
	}
	logger.Info("nothing to process.")
	return []*gorm.Order{}
}

// commit prep means move from prep to ship
func handleCommitPrep(commited []int) []*gorm.Order {
	if len(commited) > 0 {
		logger.Info("processing prep -> ship:", commited)
		processed := []*gorm.Order{}
		for _, v := range commited {
			// get order from db
			var o gorm.Order
			err := o.PopulateByID(uint(v))
			if err != nil {
				logger.Error("error getting order from db")
				return []*gorm.Order{}
			}
			// set a new status
			err = o.SetStatusTo(gorm.OrderStatusShipping)
			if err != nil {
				logger.Error("error setting new status")
				return processed
			}
			// accumulate to processed
			processed = append(processed, &o)
		}
		return processed
	}
	logger.Info("nothing to process.")
	return []*gorm.Order{}
}

// commit ship means move from ship to delivered
func handleCommitShip(commited []int) []*gorm.Order {
	if len(commited) > 0 {
		logger.Info("processing ship -> delivered:", commited)
		processed := []*gorm.Order{}
		for _, v := range commited {
			// get order from db
			var o gorm.Order
			err := o.PopulateByID(uint(v))
			if err != nil {
				logger.Error("error getting order from db")
				return []*gorm.Order{}
			}
			// set a new status
			err = o.SetStatusTo(gorm.OrderStatusDelivered)
			if err != nil {
				logger.Error("error setting new status")
				return processed
			}
			// accumulate to processed
			processed = append(processed, &o)
		}
		return processed
	}
	logger.Info("nothing to process.")
	return []*gorm.Order{}
}
