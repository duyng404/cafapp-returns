package gorm

import (
	"cafapp-returns/logger"

	"github.com/jinzhu/gorm"
)

// Product the basic product that goes on menu and orders
type Product struct {
	gorm.Model
	Tag             string // Unique tag following our internal convention
	Name            string // Internal code names
	DisplayName     string // Full Display Name on frontend
	PriceInCents    int
	ImageURL        string
	Description     string // one-line description
	DescriptionHTML string // html formatted description
	Status          uint
}

const (
	// ProductStatusCreated created, non-active
	ProductStatusCreated = 0
	// ProductStatusOnShelf on sale
	ProductStatusOnShelf = 10
	// ProductStatusAddon is on sale but not displayed on main menu
	ProductStatusAddon = 11
	// ProductStatusDiscontinued no longer on sale
	ProductStatusDiscontinued = 90
	// ProductStatusMisc uncategorized
	ProductStatusMisc = 99
)

// Create create the object
func (p *Product) Create() error {
	return DB.Create(p).Error
}

// PopulateByID : query the db to get object by id
func (p *Product) PopulateByID(id uint) error {
	return DB.Where("id = ?", id).Last(&p).Error
}

// PopulateByIDOnShelf : query the db to get object by id, has to be on shelf
func (p *Product) PopulateByIDOnShelf(id uint) error {
	return DB.Where("id = ? AND status = ?", id, ProductStatusOnShelf).Last(&p).Error
}

// GetAllProductsOnShelf : get all products currently on sale
func GetAllProductsOnShelf() ([]Product, error) {
	var res []Product
	err := DB.Raw(`
		SELECT p.*
		FROM products p
		WHERE p.status = ?
		AND p.deleted_at IS NULL
	`, ProductStatusOnShelf).Scan(&res).Error
	return res, err
}

// IsAMeal return true if the product is in the db and ON SHELF
func IsAMeal(id uint) bool {
	var tmp []Product
	err := DB.Raw(`
		SELECT p.*
		FROM products p
		WHERE p.id = ?
		AND p.status = ?
		AND p.deleted_at IS NULL
	`, id, ProductStatusOnShelf).Scan(&tmp).Error
	if err != nil {
		logger.Error("product id", id, "is not a valid meal item", err)
		return false
	}
	if len(tmp) != 1 {
		logger.Error("product id", id, "is not a valid meal item. Returned rows are not 1.")
		return false
	}
	return true
}
