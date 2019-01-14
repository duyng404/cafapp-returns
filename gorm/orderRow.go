package gorm

import (
	"github.com/jinzhu/gorm"
)

// OrderRow a row in an order
type OrderRow struct {
	gorm.Model
	// OrderID
	ProductID       uint
	Product         Product
	Quantity        int
	SubtotalInCents int
	RowType         int
}

const (
	// RowTypeNormal chargeable normally
	RowTypeNormal = 1
	// RowTypeIncluded included with order
	RowTypeIncluded = 2
)

// Create create the object
func (o *OrderRow) Create() error {
	return DB.Create(o).Error
}
