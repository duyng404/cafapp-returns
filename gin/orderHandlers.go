package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleOrderGet(c *gin.Context) {
	// get step from sess
	// s := sessions.Default(c)
	// step := getStringFromSession(s, "currentStep")
	urlpath := c.Param("stuff")
	if urlpath == "" {
		// no existing order, start a new one
		getOrderMenu(c)
		return
	}

	// check if urlpath has a valid order uuid
	if urlpath[0] == '/' {
		urlpath = urlpath[1:]
	}
	var order gorm.Order
	err := order.PopulateByUUID(urlpath)
	if err != nil {
		logger.Error("order uuid in params is not valid:", err)
		// TODO handle error
	}

	// does order need any more info
	if order.StatusCode == 1 || order.StatusCode == 2 || order.StatusCode == 3 {
		getMoreInfo(c, order)
	}
}

func handleOrderPost(c *gin.Context) {
	// get step from sess
	// s := sessions.Default(c)
	// step := getStringFromSession(s, "currentStep")
	postOrderMenu(c)
}

// GET step 1: show the menu
func getOrderMenu(c *gin.Context) {
	// get all menu items from db
	menu, err := gorm.GetAllProductsOnShelf()
	if err != nil {
		logger.Error("could not get products to display:", err)
		// TODO handle error
	}
	// render
	renderHTML(c, 200, "order-menu.html", gin.H{
		"menu": menu,
	})
}

// GET step 2: ask the user more info to complete the order
func getMoreInfo(c *gin.Context, order gorm.Order) {
	data := make(map[string]interface{})

	// get current user
	user := getCurrentAuthUser(c)
	if user == nil {
		logger.Error("cannot get currently logged in user")
		// TODO handle error
	}

	// does user have gus id
	if user.GusID == 0 {
		data["need-gus-id"] = true
	}

	renderHTML(c, 200, "order-info.html", data)
}

// POST step 1: user has selected an item
func postOrderMenu(c *gin.Context) {
	// get current user
	user := getCurrentAuthUser(c)
	if user == nil {
		logger.Error("cannot get currently logged in user")
		// TODO handle error
	}

	// get selected
	tmp := c.PostForm("selected-meal")
	selected, err := strconv.Atoi(tmp)
	if err != nil {
		logger.Error("error getting selected meal ("+tmp+") from POST form", err)
		// TODO handle error
	}

	// get the product from db
	var p gorm.Product
	err = p.PopulateByIDOnShelf(uint(selected))
	if err != nil {
		logger.Error("error getting product from db", err)
		// TODO handle error
	}

	// create a new order
	row := gorm.NewOrderRowFromProduct(&p)
	order := gorm.Order{
		User:       user,
		OrderRows:  []gorm.OrderRow{*row},
		StatusCode: 1,
	}
	err = order.Create()
	if err != nil {
		logger.Error("error creating new order", err)
		// TODO handle error
	}

	logger.Info("created order with uuid", order.UUID)

	c.Redirect(http.StatusFound, "/order/"+order.UUID)
}
