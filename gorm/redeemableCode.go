package gorm

import (
	"errors"
	"time"

	"github.com/AmirSoleimani/VoucherCodeGenerator/vcgen"
	"github.com/jinzhu/gorm"
)

// RedeemableCode are the codes that can be redeemed to add Delivery Balance. These codes goes on sale at the Bookmark.
type RedeemableCode struct {
	Code             string     `json:"code" gorm:"primary_key"`
	CreatedAt        time.Time  `json:"created_at"`
	AmountInCents    int        `json:"amount_in_cents"`
	Status           int        `json:"status"`
	RedeemedAt       *time.Time `json:"redeemed_at"`
	RedeemedByUserID uint       `json:"redeemed_by_user_id"`
	RedeemedByUser   User       `json:"redeemed_by_user" gorm:"foreignkey:RedeemedByUserID,association_foreignkey:ID"`
}

const (
	// RedeemableCodeStatusAvailable available for redeem
	RedeemableCodeStatusAvailable = 1
	// RedeemableCodeStatusOnSale available for redeem and has been put on sale
	RedeemableCodeStatusOnSale = 2
	// RedeemableCodeStatusRedeemed already been redeemed
	RedeemableCodeStatusRedeemed = 3
)

// GetRedeemableCodeByCode input a string, will query db and return a RedeemableCode object with that string as the Code
func GetRedeemableCodeByCode(code string) (*RedeemableCode, error) {
	var c RedeemableCode
	err := DB.Where("code = ?", code).Last(&c).Error
	return &c, err
}

// GetAllRedeemableCodes for viewing in admin dash
func GetAllRedeemableCodes() ([]RedeemableCode, error) {
	var res []RedeemableCode
	err := DB.Preload("RedeemedByUser").Find(&res).Error
	return res, err
}

// GenerateFiveRedeemableCodes : generate a set of 5 codes and save it to db
func GenerateFiveRedeemableCodes() ([]*RedeemableCode, error) {
	var res []*RedeemableCode
	// loop until we have 5
	for len(res) < 5 {
		// generate
		generator := vcgen.New(&vcgen.Generator{
			Count:   5,
			Pattern: "####-####-####",
			Prefix:  "CA-",
			Charset: "123456789QWERTYUPADFGHJKLXCVBNM",
		})
		result, _ := generator.Run()
		for _, v := range *result {
			// check if already in DB
			_, err := GetRedeemableCodeByCode(v)
			if err == nil {
				continue
			} else if err != gorm.ErrRecordNotFound {
				return res, err
			}
			// create object in db
			code := RedeemableCode{
				Code:          v,
				AmountInCents: 1000,
				Status:        RedeemableCodeStatusAvailable,
			}
			err = DB.Create(&code).Error
			if err != nil {
				return res, err
			}
			// accumulate
			res = append(res, &code)
		}
	}
	return res, nil
}

// RedeemACodeByCode set the code status to redeemed, and add its value to user's balance
func RedeemACodeByCode(user User, code string) (bool, error) {
	if user.ID == 0 {
		return false, errors.New("id is zero")
	}
	c, err := GetRedeemableCodeByCode(code)
	if err != nil {
		return false, err
	}
	if c.Status != RedeemableCodeStatusAvailable {
		return false, nil
	}
	c.Status = RedeemableCodeStatusRedeemed
	c.RedeemedByUser = user
	now := time.Now()
	c.RedeemedAt = &now
	// TODO: increase the user's balance somewhere.
	err = DB.Save(c).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
