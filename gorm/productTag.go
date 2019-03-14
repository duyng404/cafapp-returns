package gorm

import (
	"github.com/jinzhu/gorm"
)

// Tag is a way to categorize products (into meals, drinks, sides etc)
type Tag struct {
	gorm.Model
	Name string
}

// GetOrCreateTag retrieve tag from db. if not exist, create, then return
func GetOrCreateTag(name string) (*Tag, error) {
	var t Tag
	err := DB.Where("name = ?", name).Last(&t).Error
	if err == ErrRecordNotFound {
		t.Name = name
		err = DB.Create(&t).Error
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func initTags() {
	GetOrCreateTag("main")
	GetOrCreateTag("side")
	GetOrCreateTag("drink")
}
