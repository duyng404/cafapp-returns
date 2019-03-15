package gorm

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"
)

// Order : the most important object in our app
type Order struct {
	gorm.Model
	UUID                          string              `json:"uuid" gorm:"index:uuid"`
	Tag                           string              `json:"tag"`
	UserID                        uint                `json:"user_id"`
	User                          *User               `json:"user"`
	DeliveryFeeInCents            int                 `json:"delivery_fee_in_cents"`
	CafAccountChargeAmountInCents int                 `json:"caf_account_charge_amount_in_cents"`
	TotalInCents                  int                 `json:"total_in_cents"`
	OrderRows                     []OrderRow          `json:"order_rows" gorm:"many2many:order_order_rows"`
	DestinationTag                string              `json:"destination_tag"`
	Destination                   *Destination        `json:"destination" gorm:"foreignkey:DestinationTag,association_foreignkey:Tag"`
	StatusCode                    int                 `json:"status_code"`
	StatusUpdates                 []OrderStatusUpdate `jons:"order_status_updates"`
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
	return DB.
		Preload("User").
		Preload("OrderRows", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_rows.id") // Preload OrderRows and sort them by order_rows.id
		}).
		Preload("OrderRows.MenuItem").
		Preload("OrderRows.SubRows", func(db *gorm.DB) *gorm.DB {
			return db.Order("sub_rows.id") // Preload OrderRows.SubRows and sort them by sub_rows.id
		}).
		Preload("OrderRows.SubRows.Product").
		Preload("Destination").
		Preload("StatusUpdates").
		Where("id = ?", id).Last(&o).Error
}

// PopulateByUUID : query the db to get object by uuidid
func (o *Order) PopulateByUUID(uuid string) error {
	return DB.
		Preload("User").
		Preload("OrderRows", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_rows.id") // Preload OrderRows and sort them by order_rows.id
		}).
		Preload("OrderRows.MenuItem").
		Preload("OrderRows.SubRows", func(db *gorm.DB) *gorm.DB {
			return db.Order("sub_rows.id") // Preload OrderRows.SubRows and sort them by sub_rows.id
		}).
		Preload("OrderRows.SubRows.Product").
		Preload("Destination").
		Preload("StatusUpdates").
		Where("uuid = ?", uuid).Last(&o).Error
}

// GetDeliveredTime : return the time when the order was delivered
func GetDeliveredTime(id uint) (OrderStatusUpdate, error) {
	var tmp OrderStatusUpdate
	err := DB.Raw(`
		SELECT order_status_updates.*
		FROM order_status_updates
		WHERE status_code = ? AND order_id = ?
	`, OrderStatusDelivered, id).Scan(&tmp).Error
	if tmp.ID == 0 {
		return tmp, ErrRecordNotFound
	}
	return tmp, err
}

// GetAllOrderFromUser : return all Orders that Users have placed
func GetAllOrderFromUser(id uint) (*[]Order, error) {
	var orders []Order
	err := DB.Preload("User").Preload("OrderRows").Preload("OrderRows").Preload("Destination").
		Where("user_id = ? AND status_code >= ?", id, OrderStatusPlaced).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

// CalculateDeliveryFee : calculate the delivery of a given order
// does not save. Caller should handle that
func (o *Order) CalculateDeliveryFee() {
	o.DeliveryFeeInCents = 250
}

// CalculateTotal : calculate the total fee based on what's in order rows
// does not save. Caller should handle that
func (o *Order) CalculateTotal() {
	total := 0
	for i := range o.OrderRows {
		subtotal := 0
		for _, v := range o.OrderRows[i].SubRows {
			if v.Product != nil {
				subtotal += v.Product.PriceInCents
			}
		}
		o.OrderRows[i].SubtotalInCents = subtotal
		total += subtotal
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
	orderNum, err := g.GetNextOrderTag()
	if err != nil {
		return err
	}
	tagString := ""
	for i := range o.OrderRows {
		for j := range o.OrderRows[i].SubRows {
			tagString = tagString + o.OrderRows[i].SubRows[j].Product.Tag
		}
	}
	o.Tag = fmt.Sprintf("%s-%s-%d", o.DestinationTag, tagString, orderNum)
	return nil
}

// GetOrdersForAdminViewQueue :
func GetOrdersForAdminViewQueue() (*[]Order, error) {
	var orders []Order
	twentyFourHours := time.Duration(24) * time.Hour
	twentyFourHoursFromNow := time.Now().Add(-twentyFourHours)
	err := DB.
		Order("tag").
		Preload("User").
		Preload("OrderRows").
		Preload("OrderRows.Product").
		Preload("Destination").
		Where("status_code >= ? AND status_code < ? AND created_at > ?", OrderStatusPlaced, OrderStatusDelivered, twentyFourHoursFromNow).Find(&orders).Error
	return &orders, err
}

// SetStatusTo : changes the status and save the update to the db.
func (o *Order) SetStatusTo(s int) error {
	if o.ID == 0 {
		return errors.New("id is zero")
	}
	o.StatusCode = s
	err := CreateOrderStatusUpdate(o.ID, s)
	if err != nil {
		return err
	}
	err = o.Save()
	if err != nil {
		return err
	}
	return nil
}
