package gin

import (
	"github.com/gin-gonic/gin"
)

func handleLandingTop(c *gin.Context) {
	renderHTML(c, "landing-top.html", gin.H{})
}

func handleLandingAbout(c *gin.Context) {
	renderHTML(c, "landing-about.html", gin.H{})
}

func handleLandingNews(c *gin.Context) {
	renderHTML(c, "landing-news.html", gin.H{})
}

func handleLandingMenu(c *gin.Context) {
	renderHTML(c, "landing-menu.html", gin.H{})
}
