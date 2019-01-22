package gin

import (
	"cafapp-returns/config"
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
