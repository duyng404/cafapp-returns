package gorm

import (
	//"cafapp-returns/logger"
	"errors"
	"fmt"

	//"github.com/davecgh/go-spew/spew"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"
)

// Order : the most important object in our app
type Order struct {
	gorm.Model
	UUID                          string       `json:"uuid" gorm:"index:uuid"`
	Tag                           string       `json:"tag"`
	UserID                        uint         `json:"user_id"`
	User                          *User        `json:"user"`
	DeliveryFeeInCents            int          `json:"delivery_fee_in_cents"`
	CafAccountChargeAmountInCents int          `json:"caf_account_charge_amount_in_cents"`
	TotalInCents                  int          `json:"total_in_cents"`
	OrderRows                     []OrderRow   `json:"order_rows" gorm:"many2many:order_order_rows"`
	DestinationTag                string       `json:"destination_tag"`
	Destination                   *Destination `json:"destination" gorm:"foreignkey:DestinationTag,association_foreignkey:Tag"`
	StatusCode                    int          `json:"status_code"`
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
	return DB.Preload("User").
		Preload("OrderRows").
		Preload("OrderRows.Product").
		Preload("Destination").
		Where("id = ?", id).Last(&o).Error
}

// PopulateByUUID : query the db to get object by uuidid
func (o *Order) PopulateByUUID(uuid string) error {
	return DB.
		Preload("User").
		Preload("OrderRows").
		Preload("OrderRows.Product").
		Preload("Destination").
		Where("uuid = ?", uuid).Last(&o).Error
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

// GetAllOrderFromUser : return all Orders that Users have placed
func GetAllOrderFromUser(id uint) (*[]Order, error) {
	var orders []Order
	err := DB.Preload("User").Preload("OrderRows").Preload("OrderRows.Product").
		Where("user_id = ?", id).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	//logger.Info("data is: ", spew.Sdump(orders))
	return &orders, nil
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

// GenerateTag :
func (o *Order) GenerateTag() error {
	g, err := GetGlobalVar()
	if err != nil {
		return err
	}
	tag, err := g.GetNextOrderTag()
	if err != nil {
		return err
	}
	meal := o.GetMealRow()
	drink := o.GetDrinkRow()
	if meal == nil || drink == nil {
		return errors.New("order is not completed")
	}
	o.Tag = fmt.Sprintf("%s-%s%s-%d", o.DestinationTag, meal.Product.Tag, drink.Product.Tag, tag)
	return nil
}

// GetOrdersForAdminViewQueue :
func GetOrdersForAdminViewQueue() (*[]Order, error) {
	var orders []Order
	err := DB.
		Preload("User").
		Preload("OrderRows").
		Preload("OrderRows.Product").
		Preload("Destination").
		Where("status_code BETWEEN ? AND ?", OrderStatusPlaced, OrderStatusDelivered).Find(&orders).Error
	return &orders, err
}
