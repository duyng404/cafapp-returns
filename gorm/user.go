package gorm

import (
	"cafapp-returns/apiObjects"
	"cafapp-returns/jwt"
	"cafapp-returns/logger"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// User : a cafapp user!
type User struct {
	gorm.Model
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	FullName              string `json:"full_name"`
	Email                 string `json:"email" gorm:"index:email"`
	GusUsername           string `json:"gus_username"`
	GusID                 int    `json:"gus_id"`
	IsAdmin               bool   `json:"-"`
	CurrentBalanceInCents int    `json:"current_balance_in_cents"`
	PhoneNumber           string `json:"phone_number"`
}

// GetAllUser ...
func (u *User) GetAllUser() ([]User, error) {
	var users []User
	err := DB.Raw(`
		SELECT u.*
		FROM users u
		WHERE u.deleted_at IS NULL
	`).Scan(&users).Error
	return users, err
}

// Create : Create the object
func (u *User) Create() error {
	return DB.Create(u).Error
}

// Save : save the object
func (u *User) Save() error {
	return DB.Save(u).Error
}

// PopulateByID : query the db to get object by id
func (u *User) PopulateByID(id uint) error {
	return DB.Where("id = ?", id).Last(&u).Error
}

// PopulateByEmail : query the db to get object by email
func (u *User) PopulateByEmail(e string) error {
	return DB.Where("email = ?", e).Last(&u).Error
}

// GenerateJWT generates a new jwt for the user
func (u *User) GenerateJWT() (string, error) {
	// expire the token in a week
	expire := time.Now().Add(time.Hour * 168).Unix()
	// return jwt.NewToken(u.ID, u.GusUsername, u.IsAdmin, expire)
	token, err := jwt.NewToken(u.ID, u.GusUsername, u.IsAdmin, expire)
	if err != nil {
		logger.Error(err)
	}
	return token, err
}

// GenerateSocketToken generate a new socket token for the user
func (u *User) GenerateSocketToken() (string, error) {
	token := UserSocketToken{
		User:   *u,
		UserID: u.ID,
	}
	err := token.Renew()
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

// GetOneIncompleteOrder : check the latest order. If it has status <= finalized, return it.
func (u *User) GetOneIncompleteOrder() (*Order, error) {
	var o Order
	err := DB.Raw(`
		SELECT o.*
		FROM orders o
		WHERE o.user_id = ?
		ORDER BY o.created_at DESC
		LIMIT 1
	`, u.ID).Scan(&o).Error
	if err != nil {
		return nil, err
	}
	if o.StatusCode > OrderStatusFinalized {
		return nil, nil
	}
	return &o, nil
}

// RedeemDeliveryCode : add the delivery code's amount to user's balance, and create a transaction record
func (u *User) RedeemDeliveryCode(code *RedeemableCode) (bool, error) {
	// mark the code as redeemed
	ok, err := code.MarkCodeAsRedeemed(u)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	// increase balance
	originalAmount := u.CurrentBalanceInCents
	if ok {
		u.CurrentBalanceInCents += code.AmountInCents
		err := u.Save()
		if err != nil {
			return false, err
		}
	}

	// save transaction
	t := Transaction{
		User:         u,
		AmountBefore: originalAmount,
		AmountAfter:  u.CurrentBalanceInCents,
		Type:         TransactionTypeRedeem,
		RelatedInfo:  code.Code,
	}
	err = t.Create()
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetUsersForAdmin get all users for admin dashboard, with filters and sorting
func GetUsersForAdmin(fn string, gususername string, sortBy string) ([]apiObjects.AdminUsersStruct, error) {
	var tmp []apiObjects.AdminUsersStruct
	var sql strings.Builder
	switch sortBy {
	case "idDESC":
		sortBy = "u.id DESC"
	case "full_nameDESC":
		sortBy = "u.full_name DESC"
	case "gus_usernameDESC":
		sortBy = "u.gus_username DESC"
	}
	sql.WriteString(`SELECT
			u.*,
			(SELECT COUNT(*)
			FROM orders o_sub
			WHERE o_sub.user_id = u.id
			AND o_sub.status_code >= ?) AS total_orders
			FROM users u
			`)
	//both fullname and gususername are empty
	if len(fn) == 0 && len(gususername) == 0 {
		err := DB.Raw(sql.String(), OrderStatusPlaced).Order(sortBy).Scan(&tmp).Error
		return tmp, err
	} else if len(fn) > 0 && len(gususername) == 0 {
		sql.WriteString(`WHERE full_name LIKE ?`)
		err := DB.Raw(sql.String(), OrderStatusPlaced, "%"+fn+"%").Order(sortBy).Scan(&tmp).Error
		return tmp, err
	} else if len(fn) == 0 && len(gususername) > 0 {
		sql.WriteString(`WHERE gus_username LIKE ?`)
		err := DB.Raw(sql.String(), OrderStatusPlaced, "%"+gususername+"%").Order(sortBy).Scan(&tmp).Error
		return tmp, err
	} else {
		sql.WriteString(`WHERE full_name LIKE ? AND gus_username LIKE ?`)
		err := DB.Raw(sql.String(), OrderStatusPlaced, "%"+fn+"%", "%"+gususername+"%").Order(sortBy).Scan(&tmp).Error
		return tmp, err
	}
}

//PopulateByIDForAdminDash get info for one user (admin)
func PopulateByIDForAdminDash(id uint) (apiObjects.AdminUsersStruct, error) {
	var user apiObjects.AdminUsersStruct
	err := DB.Raw(`
		SELECT
			u.*,
			(SELECT COUNT(*)
			FROM redeemable_codes r
			WHERE r.redeemed_by_user_id = u.id
			AND r.status = ?) AS number_of_redeems
		FROM users u
		WHERE u.id = ?
		AND u.deleted_at IS NULL
	`, RedeemableCodeStatusRedeemed, id).Scan(&user).Error
	return user, err
}

// NewOrderFromMenuItem ...
func (u *User) NewOrderFromMenuItem(mi *MenuItem) (*Order, error) {
	if mi.ID == 0 {
		return nil, ErrIDZero
	}

	newOrder := Order{
		User: u,
		OrderRows: []OrderRow{
			OrderRow{
				MenuItemID: mi.ID,
				MenuItem:   mi,
				SubRows: []SubRow{
					SubRow{
						Product:   mi.StartingMain,
						ProductID: mi.StartingMainID,
					},
					SubRow{
						Product:   mi.StartingSide,
						ProductID: mi.StartingSideID,
					},
				},
			},
		},
		StatusCode: OrderStatusNeedInfo,
	}
	newOrder.CalculateDeliveryFee()
	newOrder.CalculateTotal()
	err := newOrder.Create()
	if err != nil {
		logger.Warning("error creating new order", err)
		return nil, err
	}
	return &newOrder, nil
}

// GetActiveOrders ...
func (u *User) GetActiveOrders() ([]Order, error) {
	var orders []Order
	// err := DB.Raw(`
	// 	SELECT *
	// 	FROM orders o
	// 	WHERE
	// 		o.deleted_at IS NULL
	// 		AND o.user_id = ?
	// 		AND o.status_code >= ? AND o.status_code <= ?
	// 		AND o.created_at >= now() - INTERVAL 1 DAY;
	// `, u.ID, OrderStatusPlaced, OrderStatusDelivered).Scan(&tmp).Error
	err := DB.Preload("Destination").Where(`
			user_id = ?
			AND status_code >= ? AND status_code <= ?
			AND created_at >= now() - INTERVAL 1 DAY
		`, u.ID, OrderStatusPlaced, OrderStatusDelivered).Find(&orders).Error
	if err != nil {
		return orders, err
	}
	if len(orders) > 0 {
		return orders, nil
	}
	return orders, nil
}
