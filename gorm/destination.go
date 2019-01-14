package gorm

import (
	"cafapp-returns/logger"

	"github.com/jinzhu/gorm"
)

// Destination of deliveries
type Destination struct {
	gorm.Model
	Name    string
	ShortID string
}

var (
	destinationList = []Destination{
		Destination{
			Name:    "Sohre",
			ShortID: "SO",
		},
		Destination{
			Name:    "International Center",
			ShortID: "IC",
		},
		Destination{
			Name:    "Southwest",
			ShortID: "SW",
		},
	}
)

// PopulateByShortID query db by shortid
func (d *Destination) PopulateByShortID(shortid string) error {
	return DB.Where("short_id = ?", shortid).Last(&d).Error
}

// FirstOrCreate create if not exist
func (d *Destination) FirstOrCreate() error {
	var tmp Destination
	if err := tmp.PopulateByShortID(d.ShortID); err != nil {
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
