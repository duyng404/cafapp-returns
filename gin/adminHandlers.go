package gin

import (
	"cafapp-returns/config"
	"cafapp-returns/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAdminDash(c *gin.Context) {
	c.Redirect(http.StatusFound, config.AdminDashboardURL)
}

func handleAdminInfo(c *gin.Context) {
	user := getCurrentAuthUser(c)
	c.JSON(http.StatusOK, gin.H{
		"admin_name":     user.FullName,
		"admin_username": user.GusUsername,
	})
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
