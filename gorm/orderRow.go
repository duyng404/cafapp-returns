package gorm

import (
	"cafapp-returns/logger"

	"github.com/jinzhu/gorm"
)

// OrderRow a row in an order
type OrderRow struct {
	gorm.Model
	MenuItemID      uint      `json:"menu_item_id"`
	MenuItem        *MenuItem `json:"menu_item"`
	SubtotalInCents int       `json:"subtotal_in_cents"`
	SubRows         []SubRow
}

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
