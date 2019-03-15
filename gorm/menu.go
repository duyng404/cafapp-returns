package gorm

import (
	"github.com/jinzhu/gorm"
)

// Menu is a list of items that can go on sale. Only one menu can be on sale at a time.
// The active menu is set inside GlobalVar.
type Menu struct {
	gorm.Model
	Name        string     `json:"name"`         // internal name
	DisplayName string     `json:"display_name"` // name on frontend
	Description string     `json:"description"`
	MenuItems   []MenuItem `json:"menu_items"`
}

// Create creates the object in db
func (m *Menu) Create() error {
	return DB.Create(m).Error
}

// GetAllMenus ...
func GetAllMenus() ([]Menu, error) {
	var menus []Menu
	err := DB.Preload("MenuItems").Find(&menus).Error
	return menus, err
}
