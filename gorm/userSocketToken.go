package gorm

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"
)

// UserSocketToken : for identifying users when they connect to the Order Tracker or Admin Dashboard.
// the Socket.io interface doesn't have built-in authentication, so we have to use this as authentication instead.
type UserSocketToken struct {
	gorm.Model
	UserID uint
	User   User
	Token  string
}

// Renew will create a new token for user, overwriting the old one
func (t *UserSocketToken) Renew() error {
	if t.UserID == 0 {
		return ErrIDZero
	}
	err := t.PopulateByUserID(t.UserID)
	t.Token = shortuuid.New() + shortuuid.New()
	if err == ErrRecordNotFound {
		return DB.Create(t).Error
	}
	return DB.Save(t).Error
}

// PopulateByUserID Find the current token of user
func (t *UserSocketToken) PopulateByUserID(id uint) error {
	return DB.Where("user_id = ?", id).Last(t).Error
}

// ValidateAdminSocketToken : query the db for the owner of the token and check it token legit
func ValidateAdminSocketToken(token string) (*User, error) {
	var t UserSocketToken
	err := DB.Preload("User").Where("token = ?", token).Last(&t).Error
	if err != nil {
		return nil, err
	}
	if t.User.IsAdmin == false {
		return nil, errors.New("not an admin")
	}
	// token only valid for 30 minutes
	if time.Since(t.UpdatedAt).Minutes() > 30 {
		return nil, errors.New("token expired")
	}
	return &t.User, nil
}

// ValidateUserSocketToken : query the db for the owner of the token and check if token legit
func ValidateUserSocketToken(token string) (*User, error) {
	var t UserSocketToken
	err := DB.Preload("User").Where("token = ?", token).Last(&t).Error
	if err != nil {
		return nil, err
	}
	// token only valid for 30 minutes
	if time.Since(t.UpdatedAt).Minutes() > 30 {
		return nil, errors.New("token expired")
	}
	return &t.User, nil
}
