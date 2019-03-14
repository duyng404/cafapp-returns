package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleLandingTop(c *gin.Context) {
	isrunning, err := gorm.IsCafAppRunning()
	if err != nil {
		logger.Error("database error:", err)
	}
	renderHTML(c, 200, "landing-top.html", gin.H{
		"CafAppRunning": isrunning,
	})
}

func handleLandingAbout(c *gin.Context) {
	isrunning, err := gorm.IsCafAppRunning()
	if err != nil {
		logger.Error("database error:", err)
	}
	renderHTML(c, 200, "landing-about.html", gin.H{
		"Title":         "About Us",
		"CafAppRunning": isrunning,
	})
}

func handleLandingNews(c *gin.Context) {
	isrunning, err := gorm.IsCafAppRunning()
	if err != nil {
		logger.Error("database error:", err)
	}
	renderHTML(c, 200, "landing-news.html", gin.H{
		"Title":         "News",
		"CafAppRunning": isrunning,
	})
}

func handleLandingMenu(c *gin.Context) {
	isrunning, err := gorm.IsCafAppRunning()
	if err != nil {
		logger.Error("database error:", err)
	}
	user := getCurrentAuthUser(c)
	if user != nil && isrunning {
		c.Redirect(http.StatusFound, "/order")
		return
	}
	menu, err := gorm.GetAllProductsOnShelf()
	if err != nil {
		logger.Error("could not get products to display:", err)
	}
	renderHTML(c, 200, "landing-menu.html", gin.H{
		"Title":         "Menu",
		"CafAppRunning": isrunning,
		"menu":          menu,
	})
}

func handleLandingFAQ(c *gin.Context) {
	renderHTML(c, 200, "landing-faq.html", gin.H{
		"Title": "F.A.Q.",
	})
}
