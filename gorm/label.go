package gorm

import (
	"github.com/jinzhu/gorm"
)

const (
	ProductLabelMain  = "main"
	ProductLabelSide  = "side"
	ProductLabelDrink = "drink"
)

// Label is a way to categorize products (into meals, drinks, sides etc)
type Label struct {
	gorm.Model
	Name string `json:"name"`
}

// GetOrCreateLabel retrieve Label from db. if not exist, create, then return
func GetOrCreateLabel(name string) (*Label, error) {
	var label Label
	err := DB.Where("name = ?", name).Last(&label).Error
	if err == ErrRecordNotFound {
		label.Name = name
		err = DB.Create(&label).Error
		if err != nil {
			return nil, err
		}
		return &label, nil
	}
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func initLabels() {
	GetOrCreateLabel(ProductLabelMain)
	GetOrCreateLabel(ProductLabelSide)
	GetOrCreateLabel(ProductLabelDrink)
}
