package gorm

import (
	"github.com/jinzhu/gorm"
)

// MenuItem is an item on the menu. This is separate from Product.
// While Product is used for official records and calculations, MenuItem is
// only for displaying on frontend.
type MenuItem struct {
	gorm.Model
	DisplayName         string
	DisplayPriceInCents int
	ImageURL            string
	Description         string
	DescriptionHTML     string
	StartingMain        Product
	StartingMainID      uint
	StartingSide        Product
	StartingSideID      uint
	Menu                Menu
	MenuID              uint
}

// Create create the object in db
func (mi *MenuItem) Create() error {
	return DB.Create(mi).Error
}
