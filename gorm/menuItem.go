package gorm

import (
	"cafapp-returns/logger"

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
	StartingMain        *Product
	StartingMainID      uint
	StartingSide        *Product
	StartingSideID      uint
	MenuID              uint
	Menu                *Menu
}

// Create create the object in db
func (mi *MenuItem) Create() error {
	return DB.Create(mi).Error
}

// PopulateByID ...
func (mi *MenuItem) PopulateByID(id uint) error {
	return DB.Preload("StartingMain").Preload("StartingSide").Where("id = ?", id).Last(mi).Error
}

// GetActiveMenuItems ...
func GetActiveMenuItems() ([]MenuItem, error) {
	gvar, err := GetGlobalVar()
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	var res []MenuItem
	err = DB.Raw(`
		SELECT *
		FROM menu_items
		WHERE
			deleted_at IS NULL
			AND menu_id = ?
	`, gvar.ActiveMenuID).Scan(&res).Error
	return res, err
}
