package gin

import (
	"cafapp-returns/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleUserInfo(c *gin.Context) {
	user := getCurrentAuthUser(c)
	token, err := user.GenerateSocketToken()
	if err != nil {
		logger.Error("error generating token for user", user.GusUsername, ":", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	c.JSON(200, gin.H{
		"full_name":    user.FullName,
		"gus_username": user.GusUsername,
		"socket_token": token,
	})
}

func handleOrderTracker(c *gin.Context) {
	data := make(map[string]interface{})
	renderHTML(c, 200, "tracker.html", data)
}
