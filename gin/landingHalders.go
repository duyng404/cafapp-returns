package gin

import (
	"github.com/gin-gonic/gin"
)

func handleLandingTop(c *gin.Context) {
	c.HTML(200, "landing-top.html", gin.H{})
}

func handleLandingAbout(c *gin.Context) {
	c.HTML(200, "landing-about.html", gin.H{})
}
