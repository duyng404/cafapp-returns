package gorm

import (
	"cafapp-returns/ggoauth"
	"sezzle/vault/logging"

	"github.com/jinzhu/gorm"
)

// GoogleUser : google user object, built on top of OauthResponse object
type GoogleUser struct {
	gorm.Model
	ggoauth.OauthResponse
	UserID uint
}

// Create : save the object to the db
func (u *GoogleUser) Create() error {
	err := DB.Create(u).Error
	if err != nil {
		logging.Error(err)
		return err
	}

	return nil
}
