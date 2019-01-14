package gorm

import (
	"github.com/jinzhu/gorm"
)

// OrderRow a row in an order
type OrderRow struct {
	gorm.Model
	// OrderID
	ProductID       uint
	Product         *Product
	Quantity        int
	SubtotalInCents int
	RowType         int
}

const (
	// RowTypeNormal items that is to be charged normally
	RowTypeNormal = 1
	// RowTypeIncluded items that are included with the order
	RowTypeIncluded = 2
)

// Create create the object
func (o *OrderRow) Create() error {
	return DB.Create(o).Error
}

// NewOrderRowFromProduct : return a pointer to a new row created from a product
// Does not run Create(). The caller should take care of that
func NewOrderRowFromProduct(p *Product) *OrderRow {
	row := OrderRow{
		ProductID:       p.ID,
		Product:         p,
		Quantity:        1,
		SubtotalInCents: p.PriceInCents,
		RowType:         RowTypeNormal,
	}
	return &row
}
