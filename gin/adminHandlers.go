package gin

import (
	"cafapp-returns/config"
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAdminDash(c *gin.Context) {
	c.Redirect(http.StatusFound, config.AdminDashboardURL)
}

func handleAdminInfo(c *gin.Context) {
	user := getCurrentAuthUser(c)
	token, err := user.GenerateSocketToken()
	if err != nil {
		logger.Error("error generating token for user", user.GusUsername, ":", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"admin_name":     user.FullName,
		"admin_username": user.GusUsername,
		"socket_token":   token,
	})
}

func handleAdminGetDestinations(c *gin.Context) {
	dest, err := gorm.GetAllDestinations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "internal server error",
		})
	}
	c.JSON(http.StatusOK, dest)
}

func handleAdminViewQueue(c *gin.Context) {
	orders, err := gorm.GetOrdersForAdminViewQueue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "internal server error",
		})
	}
	c.JSON(http.StatusOK, orders)
}
