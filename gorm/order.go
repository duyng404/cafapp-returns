package gorm

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"
)

// Order : the most important object in our app
type Order struct {
	gorm.Model
	UUID                          string `gorm:"index:uuid"`
	Tag                           string
	UserID                        uint
	User                          *User
	DeliveryFeeInCents            int
	CafAccountChargeAmountInCents int
	TotalInCents                  int
	OrderRows                     []OrderRow `gorm:"many2many:order_order_rows"`
	DestinationTag                string
	StatusCode                    int
}

// Create : save the object to the db
func (o *Order) Create() error {
	if o.UUID == "" {
		o.UUID = shortuuid.New()
	}
	return DB.Create(o).Error
}

// Save : Save / Update
func (o *Order) Save() error {
	if o.ID == 0 {
		return errors.New("id is zero")
	}
	return DB.Save(o).Error
}

// PopulateByID : query the db to get object by id
func (o *Order) PopulateByID(id uint) error {
	return DB.Preload("User").Preload("OrderRows").Preload("OrderRows.Product").Where("id = ?", id).Last(&o).Error
}

// PopulateByUUID : query the db to get object by uuidid
func (o *Order) PopulateByUUID(uuid string) error {
	return DB.Preload("User").Preload("OrderRows").Preload("OrderRows.Product").Where("uuid = ?", uuid).Last(&o).Error
}

// GetMealRow : return the OrderRow that is the meal part of the order
func (o *Order) GetMealRow() *OrderRow {
	for i := range o.OrderRows {
		if o.OrderRows[i].Product.Status == ProductStatusOnShelf {
			return &o.OrderRows[i]
		}
	}
	return nil
}

// GetDrinkRow : return the OrderRow that is the drink part of the order
func (o *Order) GetDrinkRow() *OrderRow {
	for i := range o.OrderRows {
		if o.OrderRows[i].Product.Status == ProductStatusAddon {
			return &o.OrderRows[i]
		}
	}
	return nil
}

// CalculateDeliveryFee : calculate the delivery of a given order
// does not save. Caller should handle that
func (o *Order) CalculateDeliveryFee() {
	// TODO: implement a proper rate
	o.DeliveryFeeInCents = 75
}

// CalculateTotal : calculate the total fee based on what's in order rows
// does not save. Caller should handle that
func (o *Order) CalculateTotal() {
	total := 0
	for _, v := range o.OrderRows {
		total += v.SubtotalInCents
	}
	o.CafAccountChargeAmountInCents = total
	total += o.DeliveryFeeInCents
	o.TotalInCents = total
}
