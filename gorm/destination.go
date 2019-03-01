package gorm

import (
	"cafapp-returns/logger"

	"github.com/jinzhu/gorm"
)

// Destination of deliveries
type Destination struct {
	gorm.Model
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

var (
	destinationList = []Destination{
		Destination{
			Name: "Norelius",
			Tag:  "NR",
		},
		Destination{
			Name: "Complex",
			Tag:  "CX",
		},
		Destination{
			Name: "Rundstrom",
			Tag:  "RU",
		},
		Destination{
			Name: "Uhler",
			Tag:  "UH",
		},
		Destination{
			Name: "Pittman",
			Tag:  "PM",
		},
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

// GetAllDestinations : get all destinations
func GetAllDestinations() ([]Destination, error) {
	var res []Destination
	err := DB.Raw(`
		SELECT d.*
		FROM destinations d
		WHERE d.deleted_at IS NULL
	`).Scan(&res).Error
	return res, err
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
