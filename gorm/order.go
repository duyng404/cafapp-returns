package gorm

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"
)

// Order : the most important object in our app
type Order struct {
	gorm.Model
	UUID            string
	PlacedAt        *time.Time
	ShippedAt       *time.Time
	SettledAt       *time.Time
	ShortID         string
	UserID          uint
	User            *User
	TotalInCents    int
	OrderRows       []OrderRow `gorm:"many2many:order_order_rows"`
	DestinationCode string
	StatusCode      int
}

// Create : save the object to the db
func (o *Order) Create() error {
	if o.UUID == "" {
		o.UUID = shortuuid.New()
	}
	return DB.Create(o).Error
}

// PopulateByID : query the db to get object by id
func (o *Order) PopulateByID(id uint) error {
	return DB.Where("id = ?", id).Last(&o).Error
}

// PopulateByUUID : query the db to get object by uuidid
func (o *Order) PopulateByUUID(uuid string) error {
	return DB.Where("uuid = ?", uuid).Last(&o).Error
}

// NewOrderFromOrderRow : return a pointer to a new order created from a row
// Does not run Create(). The caller should take care of that
func NewOrderFromOrderRow(r *OrderRow) *Order {
	order := Order{
		OrderRows: []OrderRow{*r},
	}
	return &order
}
