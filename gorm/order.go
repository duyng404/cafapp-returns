package gorm

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Order : the most important object in our app
type Order struct {
	gorm.Model
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
	return DB.Create(o).Error
}

// PopulateByID : query the db to get object by id
func (o *Order) PopulateByID(id uint) error {
	return DB.Where("id = ?", id).Last(&o).Error
}

// NewOrderFromOrderRow : return a pointer to a new order created from a row
// Does not run Create(). The caller should take care of that
func NewOrderFromOrderRow(r *OrderRow) *Order {
	order := Order{
		OrderRows: []OrderRow{*r},
	}
	return &order
}
