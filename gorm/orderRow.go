package gorm

import (
	"cafapp-returns/logger"

	"github.com/jinzhu/gorm"
)

// OrderRow a row in an order
type OrderRow struct {
	gorm.Model
	ProductID       uint      `json:"product_id"`
	Product         *Product  `json:"product"`
	MenuItemID      uint      `json:"menu_item_id"`
	MenuItem        *MenuItem `json:"menu_item"`
	Quantity        int       `json:"quantity"`
	SubtotalInCents int       `json:"subtotal_in_cents"`
	RowType         string    `json:"row_type"`
	SubRows         []SubRow
}

const (
	// RowTypeNormal DEPRECATED items that is to be charged normally
	RowTypeNormal = "1"
	// RowTypeIncluded DEPRECATED items that are included with the order
	RowTypeIncluded = "2"
	// RowTypeFirstCombo ..
	RowTypeFirstCombo = "meal1"
	// RowTypeSecondCombo ..
	RowTypeSecondCombo = "meal2"
)

// Create create the object
func (or *OrderRow) Create() error {
	return DB.Create(or).Error
}

// Save save the object
func (or *OrderRow) Save() error {
	return DB.Save(or).Error
}

// Delete ..
func (or *OrderRow) Delete() error {
	return DB.Delete(or).Error
}

// PopulateByID : query the db by id
func (or *OrderRow) PopulateByID(id uint) error {
	return DB.Where("id = ?", id).Scan(&or).Error
}

// NewOrderRowFromProduct : DEPRECATED return a pointer to a new row created from a product
// Does not run Create(). The caller should take care of that
func NewOrderRowFromProduct(p *Product, rowType string) *OrderRow {
	row := OrderRow{
		ProductID:       p.ID,
		Product:         p,
		Quantity:        1,
		SubtotalInCents: p.PriceInCents,
		RowType:         rowType,
	}
	return &row
}

// SetMainSubRowTo ..
func (or *OrderRow) SetMainSubRowTo(p *Product) error {
	for i := range or.SubRows {
		if or.SubRows[i].Product.IsMain() {
			or.SubRows[i].Product = p
			or.SubRows[i].ProductID = p.ID
			return or.SubRows[i].Save()
		}
	}
	newSubRow := SubRow{
		Product:    p,
		ProductID:  p.ID,
		OrderRowID: or.ID,
	}
	err := newSubRow.Create()
	if err != nil {
		logger.Warning(err)
		return err
	}
	or.SubRows = append(or.SubRows, newSubRow)
	return or.Save()
}

// SetSideSubRowTo ..
func (or *OrderRow) SetSideSubRowTo(p *Product) error {
	for i := range or.SubRows {
		if or.SubRows[i].Product.IsSide() {
			or.SubRows[i].Product = p
			or.SubRows[i].ProductID = p.ID
			return or.SubRows[i].Save()
		}
	}
	newSubRow := SubRow{
		Product:    p,
		ProductID:  p.ID,
		OrderRowID: or.ID,
	}
	err := newSubRow.Create()
	if err != nil {
		logger.Warning(err)
		return err
	}
	or.SubRows = append(or.SubRows, newSubRow)
	return or.Save()
}

// SetDrinkSubRowTo ..
func (or *OrderRow) SetDrinkSubRowTo(p *Product) error {
	for i := range or.SubRows {
		if or.SubRows[i].Product.IsDrink() {
			or.SubRows[i].Product = p
			or.SubRows[i].ProductID = p.ID
			return or.SubRows[i].Save()
		}
	}
	newSubRow := SubRow{
		Product:    p,
		ProductID:  p.ID,
		OrderRowID: or.ID,
	}
	err := newSubRow.Create()
	if err != nil {
		logger.Warning(err)
		return err
	}
	or.SubRows = append(or.SubRows, newSubRow)
	return or.Save()
}
