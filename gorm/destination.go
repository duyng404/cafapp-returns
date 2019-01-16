package gorm

import (
	"cafapp-returns/logger"

	"github.com/jinzhu/gorm"
)

// Destination of deliveries
type Destination struct {
	gorm.Model
	Name string
	Tag  string
}

var (
	destinationList = []Destination{
		Destination{
			Name: "Sohre",
			Tag:  "SO",
		},
		Destination{
			Name: "International Center",
			Tag:  "IC",
		},
		Destination{
			Name: "Southwest",
			Tag:  "SW",
		},
	}
)

// PopulateByTag query db by tag
func (d *Destination) PopulateByTag(tag string) error {
	return DB.Where("tag = ?", tag).Last(&d).Error
}

// FirstOrCreate create if not exist
func (d *Destination) FirstOrCreate() error {
	var tmp Destination
	if err := tmp.PopulateByTag(d.Tag); err == gorm.ErrRecordNotFound {
		return DB.Create(&d).Error
	}
	return nil
}

// create all destinations
func initDestinations() error {
	for i := range destinationList {
		err := destinationList[i].FirstOrCreate()
		if err != nil {
			logger.Error(err)
		}
	}
	return nil
}
