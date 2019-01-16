package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleOrderGet(c *gin.Context) {
	// get the order id from params
	urlpath := c.Param("stuff")
	if urlpath == "" || urlpath == "/" {
		// no existing order, start a new one
		getOrderMenu(c)
		return
	}

	// check if valid order uuid
	if urlpath[0] == '/' {
		urlpath = urlpath[1:]
	}
	var order gorm.Order
	err := order.PopulateByUUID(urlpath)
	if err != nil {
		logger.Error("order uuid in params is not valid:", err)
		logger.Info("boucing back to /order")
		c.Redirect(http.StatusFound, "/order")
	}

	// does order need any more info
	if order.StatusCode == 1 || order.StatusCode == 2 || order.StatusCode == 3 {
		getMoreInfo(c, order)
	}
}

func handleOrderPost(c *gin.Context) {
	urlpath := c.Param("stuff")
	if urlpath == "" || urlpath == "/" {
		// no order id, so user wanted to start a new one
		postOrderMenu(c)
		return
	}

	// check if valid order uuid
	if urlpath[0] == '/' {
		urlpath = urlpath[1:]
	}
	var order gorm.Order
	err := order.PopulateByUUID(urlpath)
	if err != nil {
		logger.Error("order uuid in params is not valid:", err)
		logger.Info("boucing back to /order")
		c.Redirect(http.StatusFound, "/order")
	}

	postOrderInfo(c, order)
}

// will show the user the error text and a link to start over
func orderError(c *gin.Context, err string) {
	renderHTML(c, 404, "order-error.html", gin.H{
		"error": err,
	})
}

// GET step 1: show the menu
func getOrderMenu(c *gin.Context) {
	// get all menu items from db
	menu, err := gorm.GetAllProductsOnShelf()
	if err != nil {
		logger.Error("could not get products to display:", err)
		orderError(c, "Could not load menu items")
	}
	// render
	renderHTML(c, 200, "order-menu.html", gin.H{
		"menu":  menu,
		"Title": "Build Your Order",
	})
}

// GET step 2: ask the user more info to complete the order
func getMoreInfo(c *gin.Context, order gorm.Order) {
	data := make(map[string]interface{})
	data["Title"] = "Build Your Order"

	// get current user
	user := getCurrentAuthUser(c)
	if user == nil {
		logger.Error("cannot get currently logged in user")
		orderError(c, "Authentication Error")
	}

	// does user have gus id
	if user.GusID == 0 {
		data["needGusID"] = true
	}

	// determine currently selected meal and drink
	for i := range order.OrderRows {
		if order.OrderRows[i].RowType == gorm.RowTypeNormal {
			data["selectedMealID"] = order.OrderRows[i].ProductID
			data["selectedMealName"] = order.OrderRows[i].Product.DisplayName
			data["selectedMealPrice"] = order.OrderRows[i].Product.PriceInCents
		}
		if order.OrderRows[i].RowType == gorm.RowTypeIncluded {
			data["selectedDrinkID"] = order.OrderRows[i].ProductID
			data["selectedDrinkName"] = order.OrderRows[i].ProductID
		}
	}

	// determine currently selected destination
	if order.DestinationTag != "" {
		data["selectedDestination"] = order.DestinationTag
	}

	// pass in the order details
	data["deliveryFee"] = order.DeliveryFeeInCents
	data["orderTotal"] = order.TotalInCents
	data["cafAccountChargeAmount"] = order.CafAccountChargeAmountInCents

	// get all menu items from db
	menu, err := gorm.GetAllProductsOnShelf()
	if err != nil {
		logger.Error("could not get products to display:", err)
		orderError(c, "Database Error")
	}
	data["menu"] = menu

	// get all drinks
	drinks, err := gorm.GetAllAddonProducts()
	if err != nil {
		logger.Error("could not get drinks to display:", err)
		orderError(c, "Database Error")
	}
	data["drinks"] = drinks

	// get all destinations
	dest, err := gorm.GetAllDestinations()
	if err != nil {
		logger.Error("could not load destinations:", err)
		orderError(c, "Database Error")
	}
	data["destinations"] = dest

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
		StatusCode: gorm.OrderStatusNeedInfo,
	}
	order.CalculateDeliveryFee()
	order.CalculateTotal()
	err = order.Create()
	if err != nil {
		logger.Error("error creating new order", err)
		// TODO handle error
	}

	logger.Info("created order with uuid", order.UUID)

	c.Redirect(http.StatusFound, "/order/"+order.UUID)
}

// POST step 2: user finished filling out info
func postOrderInfo(c *gin.Context, order gorm.Order) {
	logger.Info("meal is", c.PostForm("meal"))
	logger.Info("drink is", c.PostForm("drink"))
	logger.Info("dest is", c.PostForm("destination"))
	logger.Info("gusid is", c.PostForm("gus-id"))
	c.JSON(200, order)
}

// when frontend send an ajax request requesting the new price, so the price can
// be updated without reloading the page
func handleRecalculateOrder(c *gin.Context) {
	type reqStruct struct {
		MealID  uint `json:"meal_id"`
		DrinkID uint `json:"drink_id"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// build a temporary order just to calculate the price
	order := gorm.Order{}
	meal := gorm.Product{}
	drink := gorm.Product{}
	if req.MealID != 0 {
		err = meal.PopulateByID(req.MealID)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		mealrow := gorm.NewOrderRowFromProduct(&meal)
		order.OrderRows = append(order.OrderRows, *mealrow)
	}
	if req.DrinkID != 0 {
		err = drink.PopulateByID(req.DrinkID)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		drinkrow := gorm.NewOrderRowFromProduct(&drink)
		order.OrderRows = append(order.OrderRows, *drinkrow)
	}
	order.CalculateDeliveryFee()
	order.CalculateTotal()

	// return
	type resStruct struct {
		MealName         string `json:"meal_name"`
		MealPrice        string `json:"meal_price"`
		DrinkName        string `json:"drink_name"`
		DeliveryFee      string `json:"delivery_fee"`
		OrderTotal       string `json:"order_total"`
		CafAcctChargeAmt string `json:"caf_acct_charge_amt"`
	}
	res := resStruct{
		MealName:         meal.DisplayName,
		MealPrice:        formatMoney(meal.PriceInCents),
		DrinkName:        drink.DisplayName,
		DeliveryFee:      formatMoney(order.DeliveryFeeInCents),
		OrderTotal:       formatMoney(order.TotalInCents),
		CafAcctChargeAmt: formatMoney(order.CafAccountChargeAmountInCents),
	}

	c.JSON(200, res)
}
