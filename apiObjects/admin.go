package apiObjects

type AdminUsersStruct struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	FullName    string `json:"full_name"`
	Email       string `json:"email" gorm:"index:email"`
	GusUsername string `json:"gus_username"`
	GusID       int    `json:"gus_id"`
	TotalOrders int    `json:"total_orders"`
}
