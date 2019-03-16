package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"fmt"
	"net/http"
	"time"

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
	menu, err := gorm.GetActiveMenuItems()
	if err != nil {
		logger.Error("could not get products to display:", err)
	}
	renderHTML(c, 200, "landing-menu.html", gin.H{
		"Title":         "Menu",
		"CafAppRunning": isrunning,
		"menu":          menu,
	})
}

func timein(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func handleLandingFAQ(c *gin.Context) {

	for _, name := range []string{
		"",
		"Local",
		"America/Chicago",
		"Greenwich",
	} {
		t, err := timein(time.Now(), name)
		if err == nil {
			fmt.Println(t.Location(), t.Format("15:04"))
		} else {
			fmt.Println(name, "<time unknown>")
		}
	}
	renderHTML(c, 200, "landing-faq.html", gin.H{
		"Title": "F.A.Q.",
	})
}
