package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"strconv"

	"github.com/davecgh/go-spew/spew"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func handleOrderGet(c *gin.Context) {
	// get step from sess
	s := sessions.Default(c)
	step := getStringFromSession(s, "currentStep")
	switch step {
	case "1":
		logger.Info("unhandled")
	default:
		// no existing order, start a new one
		getOrderMenu(c, s)
	}
}

func handleOrderPost(c *gin.Context) {
	// get step from sess
	s := sessions.Default(c)
	step := getStringFromSession(s, "currentStep")
	// TODO also check the post form "step-number" in case user press back button
	switch step {
	case "1":
		postOrderMenu(c, s)
	default:
		// invalid
		// TODO handle error
		logger.Info("unhandled")
	}
}

// GET step 1: show the menu
func getOrderMenu(c *gin.Context, s sessions.Session) {
	// get all menu items form db
	menu, err := gorm.GetAllProductsOnShelf()
	if err != nil {
		logger.Error("could not get products to display:", err)
		// TODO handle error
	}
	// save step number
	s.Set("currentStep", "1")
	s.Save()
	// render
	renderHTML(c, 200, "order-menu.html", gin.H{
		"menu": menu,
	})
}

// POST step 1: user has selected an item
func postOrderMenu(c *gin.Context, s sessions.Session) {
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
		User:      user,
		OrderRows: []gorm.OrderRow{*row},
	}
	err = order.Create()
	if err != nil {
		logger.Error("error creating new order", err)
		// TODO handle error
	}

	logger.Info("created order is:", spew.Sdump(order))

	c.Status(200)
}
