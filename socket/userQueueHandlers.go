package socket

import (
	"cafapp-returns/gorm"
	"strconv"
)

func (c *client) handleChatbotRequest(request string, user *gorm.User) string {
	count := 0
	response := ""
	orders, err := gorm.GetAllOrderFromUser(user.ID)
	if err != nil {
		return "You haven't placed any order"
	}
	if request == "orders" {
		for i := range orders {
			if orders[i].StatusCode >= 20 && orders[i].StatusCode < 60 {
				count++
				num := strconv.Itoa(count)
				response += orders[i].Tag + ", "
				return "You have " + num + " orders in queue: " + response
			}
		}
	}
	if request == "help" {
		return "Type orders for your order's details"
	}
	return "Command not recognized. Type help for more command"
}
