package gin

import (
	"github.com/gin-gonic/gin"
)

func handleUserDash(c *gin.Context) { 
	user := getCurrentAuthUser(c)
	renderHTML(c,200, "landing-dashboard.html",gin.H{
		"Title": user.GusUsername,
		"Fullname": user.FullName,
		})
}
