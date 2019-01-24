package gin

import (
	"cafapp-returns/logger"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func handleLandingTop(c *gin.Context) {
	renderHTML(c, 200, "landing-top.html", gin.H{})
}

func handleLandingAbout(c *gin.Context) {
	renderHTML(c, 200, "landing-about.html", gin.H{
		"Title": "About Us",
	})
}

func handleLandingNews(c *gin.Context) {
	renderHTML(c, 200, "landing-news.html", gin.H{
		"Title": "News",
	})
}

func handleLandingMenu(c *gin.Context) {
	s := sessions.Default(c)

	value := s.Get("key")
	s.Delete("key")
	s.Save()
	logger.Info(value)

	renderHTML(c, 200, "landing-menu.html", gin.H{
		"Title": "Menu",
	})
}
