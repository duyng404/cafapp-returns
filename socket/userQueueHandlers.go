package socket

import (
	"cafapp-returns/gorm"
)

func (c *client) handleChatbotRequest(request string, user *gorm.User) string {
	count := 0
	response := ""
	orders, err := gorm.GetAllOrderFromUser(user.ID)
	if err != nil {
		return "You haven't placed any order"
	}
	if request == "orders"{
		for i := range orders {
			if orders[i].StatusCode >= 20 && orders[i].StatusCode < 60 {
				count++
				num := string(count)
				response += orders[i].UUID+", "
				return "You have "+num+" orders in queue: "+response
			}  
		}
	}
	if request == "help"{
		return "Type orders for your order's details"
	}
	return ""
}
