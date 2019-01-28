package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"github.com/gin-gonic/gin"
)

func handleUserDash(c *gin.Context) { 
	user := getCurrentAuthUser(c)
	orders,err := gorm.GetAllOrderFromUser(user.ID)
	if err != nil{
		logger.Error("Cannot display past orders",err)
		return
	}
	renderHTML(c,200, "landing-dashboard.html",gin.H{
		"Title": "User Dashboard",
		"Orders": orders,
		})
}
