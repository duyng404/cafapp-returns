package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleAdminViewUsers(c *gin.Context) {
	//make users array to hold the retrieved users
	fullname := c.Query("fullname")
	gususername := c.Query("gususername")
	sortby := c.Query("sortBy")
	result, err := gorm.GetUsersForAdmin(fullname, gususername, sortby)
	if err != nil {
		logger.Error("There's an error retrieving users: ", err)
		return
	}
	c.JSON(200, result)
}

func handleAdminViewOneUser(c *gin.Context) {
	var user gorm.User
	userID, err := strconv.ParseUint(c.Param("userid"), 10, 32)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	//get a list of order object for an user
	allOrders, err := gorm.GetAllOrderFromUser(uint(userID))

	err1 := user.PopulateByID(uint(userID))
	if err1 != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(200, gin.H{
		"allOrders": allOrders,
		"userInfo":  user,
	})
}
