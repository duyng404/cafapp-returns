package gorm

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// UserSocketToken : for identifying users when they connect to the Order Tracker or Admin Dashboard.
// the Socket.io interface doesn't have built-in authentication, so we have to use this as authentication instead.
type UserSocketToken struct {
	gorm.Model
	UserID uint
	User   User
	Token  string
}

// Create create the object
func (t *UserSocketToken) Create() error {
	return DB.Create(t).Error
}

// ValidateAdminSocketToken : query the db for the owner of the token, and return nil if it's an admin
func ValidateAdminSocketToken(token string) (*User, error) {
	var t UserSocketToken
	err := DB.Preload("User").Where("token = ?", token).First(&t).Error
	if err != nil {
		return nil, err
	}
	if t.User.IsAdmin == false {
		return nil, errors.New("not an admin")
	}
	return &t.User, nil
}
