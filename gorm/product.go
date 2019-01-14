package gorm

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Product the basic product that goes on menu and orders
type Product struct {
	gorm.Model
	SKU             string `json:"-" gorm:"index:sku"`
	Name            string // Full Descriptive Name
	DisplayName     string // Display Name on frontend
	PriceInCents    uint
	ImageURL        string
	Description     string // raw description
	DescriptionHTML string // html formatted description
	Status          uint
}

const (
	// ProductStatusCreated created, non-active
	ProductStatusCreated = 0
	// ProductStatusOnShelf on sale
	ProductStatusOnShelf = 10
	// ProductStatusAddon
	ProductStatusAddon = 11
	// ProductStatusDiscontinued no longer on sale
	ProductStatusDiscontinued = 90
	// ProductStatusMisc uncategorized
	ProductStatusMisc = 99
)

// Create create the object
func (p *Product) Create() error {
	if p.SKU == "" {
		return errors.New("sku is empty")
	}
	if p.ID != 0 {
		return errors.New("id is not zero")
	}
	return DB.Create(p).Error
}

// PopulateByID : query the db to get object by id
func (p *Product) PopulateByID(id uint) error {
	return DB.Where("id = ?", id).Last(&p).Error
}

// GetAllProductsOnShelf : get all products currently on sale
func GetAllProductsOnShelf() ([]Product, error) {
	var res []Product
	err := DB.Raw(`
		SELECT p.*
		FROM products p
		WHERE p.status = ?
	`, ProductStatusOnShelf).Scan(&res).Error
	return res, err
}
