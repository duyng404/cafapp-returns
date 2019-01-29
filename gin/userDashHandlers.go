package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"github.com/gin-gonic/gin"
)

func handleUserDash(c *gin.Context) { 
	user := getCurrentAuthUser(c)
	data := make(map[string]interface{})
	data["Title"] = "User Dashboard"
	
	//display all past orders
	orders,err := gorm.GetAllOrderFromUser(user.ID)
	if err != nil{
		logger.Error("Cannot display past orders",err)
		return
	}
	data["orders"] = orders

	//current user info
	data["username"] = user.GusUsername
	data["ID"] = user.GusID
	data["fullname"] = user.FullName
	data["total"] = len(*orders)
	renderHTML(c,200, "landing-dashboard.html",data)
}
