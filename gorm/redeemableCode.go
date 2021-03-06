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
	BatchNumber      int        `json:"batch_number"`
	Status           int        `json:"status"`
	Reason           string     `json:"reason"`
	RedeemedAt       *time.Time `json:"redeemed_at"`
	RedeemedByUserID uint       `json:"redeemed_by_user_id"`
	RedeemedByUser   *User      `json:"redeemed_by_user" gorm:"foreignkey:RedeemedByUserID,association_foreignkey:ID"`
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

// GenerateRedeemableCodes : generate a set amount of codes and save it to db
func GenerateRedeemableCodes(amount int, reason string, valueInCents int) ([]*RedeemableCode, error) {
	// upper limit is 50
	if amount > 50 {
		amount = 50
	}
	g, err := GetGlobalVar()
	if err != nil {
		return nil, err
	}
	batchNum, err := g.GetNextRedeemableBatchNumber()
	if err != nil {
		return nil, err
	}
	var res []*RedeemableCode
	// loop until we have enough
	for len(res) < amount {
		// generate
		generator := vcgen.New(&vcgen.Generator{
			Count:   uint16(amount),
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
				AmountInCents: valueInCents,
				BatchNumber:   batchNum,
				Status:        RedeemableCodeStatusAvailable,
				Reason:        reason,
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

// MarkCodeAsRedeemed set the code status to redeemed, and add its value to user's balance
// Return true if redeemable. False if not. And errors if any.
func (c *RedeemableCode) MarkCodeAsRedeemed(user *User) (bool, error) {
	if user.ID == 0 || c.Code == "" {
		return false, errors.New("id is zero")
	}
	if c.Status != RedeemableCodeStatusAvailable {
		return false, nil
	}
	c.Status = RedeemableCodeStatusRedeemed
	c.RedeemedByUser = user
	now := time.Now()
	c.RedeemedAt = &now
	err := DB.Save(c).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
