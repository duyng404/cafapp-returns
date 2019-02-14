package gorm

import (
	"cafapp-returns/jwt"
	"cafapp-returns/logger"
	"time"

	"github.com/lithammer/shortuuid"

	"github.com/jinzhu/gorm"
)

// User : a cafapp user!
type User struct {
	gorm.Model
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	FullName    string `json:"full_name"`
	Email       string `json:"email" gorm:"index:email"`
	GusUsername string `json:"gus_username"`
	GusID       int    `json:"gus_id"`
	IsAdmin     bool   `json:"-"`
}

// GetAllUser
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
		User:  *u,
		Token: shortuuid.New() + shortuuid.New(),
	}
	err := token.Create()
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

// GetOneIncompleteOrder check if user has any incomplete order by searching for order
// that was created less than 48h ago.
func (u *User) GetOneIncompleteOrder() (*Order, error) {
	var o Order
	err := DB.Raw(`
		SELECT o.*
		FROM orders o
		WHERE o.user_id = ?
			AND o.status_code < ?
		ORDER BY o.created_at DESC
		LIMIT 1
	`, u.ID, OrderStatusFinalized).Scan(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

//CountTotalOrders count the total orders for a particular user
func (u *User) CountTotalOrders() (int, error) {
	type totalOrder struct {
		total int
	}
	var res totalOrder
	logger.Info(u.ID)
	err := DB.Raw(`
		SELECT COUNT(*) AS total
		FROM orders o 
		WHERE o.user_id = ? 
		AND o.status_code >= ?
	`, u.ID, OrderStatusPlaced).Scan(&res).Error
	return res.total, err
}
