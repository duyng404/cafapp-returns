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
	result, err := gorm.GetUsersForAdmin()
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
