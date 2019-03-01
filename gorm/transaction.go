package gorm

import "github.com/jinzhu/gorm"

// Transaction records all transactions of users
type Transaction struct {
	gorm.Model
	User         *User  `json:"user"`
	UserID       uint   `json:"user_id"`
	AmountBefore int    `json:"amount_before"`
	AmountAfter  int    `json:"amount_after"`
	Type         string `json:"type"`
	RelatedInfo  string `json:"related_info"`
}

const (
	// TransactionTypeRedeem : when user redeem delivery cards
	TransactionTypeRedeem = "redeem"
	// TransactionTypeOrder : when user make an order
	TransactionTypeOrder = "order"
)

// Create create the object
func (t *Transaction) Create() error {
	return DB.Create(t).Error
}
