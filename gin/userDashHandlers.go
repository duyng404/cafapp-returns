package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"

	"github.com/gin-gonic/gin"
)

func handleUserDash(c *gin.Context) {
	data := make(map[string]interface{})
	data["Title"] = "My Account"

	//display all past orders
	user := getCurrentAuthUser(c)
	orders, err := gorm.GetAllOrderFromUser(user.ID)
	if err != nil {
		logger.Error("Cannot display past orders", err)
		return
	}
	data["orders"] = orders

	//current user info
	data["user"] = user
	data["totalOrders"] = len(*orders)
	renderHTML(c, 200, "userdash-top.html", data)
}
