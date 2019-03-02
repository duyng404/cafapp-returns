package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	userID, err := strconv.ParseUint(c.Param("userid"), 10, 32)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	//get a list of order object for an user
	allOrders, err := gorm.GetAllOrderFromUser(uint(userID))

	user, err := gorm.PopulateByIDForAdminDash(uint(userID))

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(200, gin.H{
		"allOrders": allOrders,
		"userInfo":  user,
	})
}
