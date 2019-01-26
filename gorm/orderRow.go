package gorm

import (
	"github.com/jinzhu/gorm"
)

// OrderRow a row in an order
type OrderRow struct {
	gorm.Model
	// OrderID
	ProductID       uint     `json:"product_id"`
	Product         *Product `json:"product"`
	Quantity        int      `json:"quantity"`
	SubtotalInCents int      `json:"subtotal_in_cents"`
	RowType         int      `json:"row_type"`
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

// Save save the object
func (o *OrderRow) Save() error {
	return DB.Save(o).Error
}

// PopulateByID : query the db by id
func (o *OrderRow) PopulateByID(id uint) error {
	return DB.Where("id = ?", id).Scan(&o).Error
}

// NewOrderRowFromProduct : return a pointer to a new row created from a product
// Does not run Create(). The caller should take care of that
func NewOrderRowFromProduct(p *Product) *OrderRow {
	var rowtype int
	if p.Status == ProductStatusOnShelf {
		rowtype = RowTypeNormal
	} else {
		rowtype = RowTypeIncluded
	}
	row := OrderRow{
		ProductID:       p.ID,
		Product:         p,
		Quantity:        1,
		SubtotalInCents: p.PriceInCents,
		RowType:         rowtype,
	}
	return &row
}
