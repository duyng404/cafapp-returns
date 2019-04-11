package gorm

import (
	"cafapp-returns/logger"
	"time"

	"github.com/jinzhu/gorm"
)

// Destination of deliveries
type Destination struct {
	DeletedAt      *time.Time `json:"deleted_at" sql:"index"`
	Name           string     `json:"name"`
	Tag            string     `json:"tag" gorm:"primary_key"`
	PickUpLocation string     `json:"pickup_location"`
}

var (
	destinationList = []Destination{
		Destination{
			Name:           "Norelius",
			Tag:            "NR",
			PickUpLocation: "E tower, by elevator",
		},
		Destination{
			Name:           "Complex",
			Tag:            "CX",
			PickUpLocation: "Sorensen door",
		},
		Destination{
			Name:           "Rundstrom",
			Tag:            "RU",
			PickUpLocation: "Parking lot",
		},
		Destination{
			Name:           "Uhler",
			Tag:            "UH",
			PickUpLocation: "Front of building",
		},
		Destination{
			Name:           "Pittman/Sohre",
			Tag:            "PM",
			PickUpLocation: "Front of building, outside",
		},
		Destination{
			Name:           "Southwest/IC",
			Tag:            "SW",
			PickUpLocation: "Front of Southwest",
		},
		Destination{
			Name:           "Prairie View",
			Tag:            "PV",
			PickUpLocation: "Door facing Bjorling parking lot",
		},
		Destination{
			Name:           "College View",
			Tag:            "CV",
			PickUpLocation: "First door",
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
		ORDER BY d.name
	`).Scan(&res).Error
	return res, err
}

// CreateOrUpdate create if not exist
func (d *Destination) CreateOrUpdate() error {
	var tmp Destination
	err := tmp.PopulateByTag(d.Tag)
	if err == gorm.ErrRecordNotFound {
		return DB.Create(d).Error
	} else if err != nil {
		return err
	} else {
		return DB.Save(d).Error
	}
}

// create all destinations
func initDestinations() error {
	for i := range destinationList {
		err := destinationList[i].CreateOrUpdate()
		if err != nil {
			logger.Error(err)
		}
	}
	return nil
}
