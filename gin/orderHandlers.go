package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"

	"github.com/gin-gonic/gin"
)

func handleOrder(c *gin.Context) {
	menu, err := gorm.GetAllProductsOnShelf()
	if err != nil {
		logger.Error("could not get products to display:", err)
	}
	for i := range menu {
		logger.Info(menu[i].DisplayName)
	}
	renderHTML(c, "order-menu.html", gin.H{
		"menu": menu,
	})
}
