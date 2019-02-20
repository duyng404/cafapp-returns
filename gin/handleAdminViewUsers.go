package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"github.com/gin-gonic/gin"
	"strconv"
)

func handleAdminViewUsers(c *gin.Context) {
	//make users array to hold the retrieved users
	result, err := gorm.GetUsersForAdmin()
	if err != nil {
		logger.Error("There's an error retrieving users: ", err)
		return
	}
	c.JSON(200, result)
}

// func handleGetTotalOrders(c *gin.Context) {
// 	var users gorm.User
// 	allUsers, err := users.GetAllUser()
// 	if err != nil {
// 		logger.Error("There's an error retrieving users: ", err)
// 		return
// 	}
// 	//create a list that stores each user's total orders
// 	var totalOrders []int
// 	for i := 0; i < len(allUsers); i++ {
// 		order, err := allUsers[i].CountTotalOrders()
// 		if err != nil {
// 			logger.Error("There's an error counting total orders", err)
// 			return
// 		}
// 		totalOrders = append(totalOrders, order)
// 	}
// 	c.JSON(200, gin.H{
// 		"totalOrders": totalOrders,
// 	})
// }
func handleGetUserAndAllOrdersFromUser(c *gin.Context) {
	var user gorm.User
	userID, err := strconv.ParseUint(c.Param("userid"), 16, 16)
	if err != nil {
		logger.Info("Cannot convert userid to uint")
		return
	}
	//get a list of order object for an user
	allOrders, err := gorm.GetAllOrderFromUser(uint(userID))
	//get an user object
	userInfo := user.PopulateByID(uint(userID))

	c.JSON(200, gin.H{
		"allOrders": allOrders,
		"userInfo":  userInfo,
	})
}
