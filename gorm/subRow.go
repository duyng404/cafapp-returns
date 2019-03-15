package gorm

import "github.com/jinzhu/gorm"

// SubRow ...
type SubRow struct {
	gorm.Model
	ProductID  uint     `json:"product_id"`
	Product    *Product `json:"product"`
	OrderRowID uint     `json:"order_row_id"`
}

// Create ...
func (sr *SubRow) Create() error {
	return DB.Create(sr).Error
}

// Save ...
func (sr *SubRow) Save() error {
	return DB.Save(sr).Error
}
