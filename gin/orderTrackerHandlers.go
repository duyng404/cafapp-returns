package gin

import "github.com/gin-gonic/gin"

func handleOrderTracker(c *gin.Context) {
	renderHTML(c, 200, "tracker.html", gin.H{})
}
