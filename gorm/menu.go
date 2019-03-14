package gorm

import (
	"github.com/jinzhu/gorm"
)

// Menu is a list of items that can go on sale. Only one menu can be on sale at a time.
// The active menu is set inside GlobalVar.
type Menu struct {
	gorm.Model
	Name        string // internal name
	DisplayName string // name on frontend
	Description string
}

// Create creates the object in db
func (m *Menu) Create() error {
	return DB.Create(m).Error
}
