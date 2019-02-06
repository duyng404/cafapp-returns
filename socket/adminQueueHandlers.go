package socket

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
)

// "commit" means "move the order up". this func handles moving the order up the queue.
func handleCommit(event string, committed []int) {
	// determine what event are we handling
	var targetStatus int
	switch event {
	case "qfeed-commit-queue":
		targetStatus = gorm.OrderStatusPrepping
	case "qfeed-commit-prep":
		targetStatus = gorm.OrderStatusShipping
	case "qfeed-commit-ship":
		targetStatus = gorm.OrderStatusDelivered
	default:
		logger.Error("no event provided")
		return
	}

	// process
	if len(committed) > 0 {
		logger.Info("processing", event, "with data", committed)
		processed := []*gorm.Order{}
		for _, v := range committed {
			// get order from db
			var o gorm.Order
			err := o.PopulateByID(uint(v))
			if err != nil {
				logger.Error("error getting order from db")
				return
			}
			// set a new status
			err = o.SetStatusTo(targetStatus)
			if err != nil {
				logger.Error("error setting new status")
				return
			}
			// accumulate to processed
			processed = append(processed, &o)
		}
		// notify all connected admins
		updateQueueForEveryone(processed)
		return
	}
	logger.Info("nothing to process.")
	return
}
