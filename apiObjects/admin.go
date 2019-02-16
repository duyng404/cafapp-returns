package apiObjects

type AdminUsersStruct struct {
	ID          int    `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	FullName    string `json:"full_name"`
	Email       string `json:"email" gorm:"index:email"`
	GusUsername string `json:"gus_username"`
	GusID       int    `json:"gus_id"`
	TotalOrders int    `json:"total_orders"`
}
