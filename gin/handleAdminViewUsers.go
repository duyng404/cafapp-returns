package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"github.com/gin-gonic/gin"
)

func handleAdminViewUsers(c *gin.Context) {
	//make users array to hold the retrieved users
	var users gorm.User
	allUsers, err := users.GetAllUser()
	if err != nil {
		logger.Error("There's an error retrieving users: ", err)
		return
	}
	c.JSON(200, gin.H{
		"users": allUsers,
	})
}
