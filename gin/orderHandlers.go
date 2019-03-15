package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"cafapp-returns/socket"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func handleOrderGet(c *gin.Context) {
	// get the order id from params
	stuff := c.Param("stuff")
	action := c.Param("action")
	if stuff == "" || stuff == "/" {
		// no existing order, start a new one
		getOrderMenu(c)
		return
	}

	// check if valid order uuid
	var order gorm.Order
	err := order.PopulateByUUID(stuff)
	if err != nil {
		logger.Error("order uuid in params is not valid:", err)
		logger.Info("boucing back to /order")
		c.Redirect(http.StatusFound, "/order")
	}

	if order.StatusCode == gorm.OrderStatusNeedInfo || order.StatusCode == gorm.OrderStatusIncomplete {
		if action == "" {
			getMoreInfo(c, order)
		} else {
			c.Redirect(http.StatusFound, "/order/"+stuff)
		}
	} else if order.StatusCode == gorm.OrderStatusFinalized {
		if action == "finalize" {
			getFinalize(c, order)
		} else {
			c.Redirect(http.StatusFound, "/order/"+stuff+"/finalize")
		}
	} else if order.StatusCode >= gorm.OrderStatusPlaced {
		if action == "completed" {
			getComplete(c, order)
		} else {
			c.Redirect(http.StatusFound, "/order/"+stuff+"/completed")
		}

	}
}

func handleOrderPost(c *gin.Context) {
	stuff := c.Param("stuff")
	action := c.Param("action")
	if stuff == "" || stuff == "/" {
		// no order id, so user wanted to start a new one
		postOrderMenu(c)
		return
	}

	// check if valid order uuid
	var order gorm.Order
	err := order.PopulateByUUID(stuff)
	if err != nil {
		logger.Error("order uuid in params is not valid:", err)
		logger.Info("boucing back to /order")
		c.Redirect(http.StatusFound, "/order")
	}

	if action != "finalize" {
		postOrderInfo(c, order)
	} else if action == "finalize" {
		postFinalize(c, order)
	}
}

// will show the user the error text and a link to start over
func orderError(c *gin.Context, err string) {
	renderHTML(c, 404, "order-error.html", gin.H{
		"error": err,
	})
}

// GET step 1: show the menu
func getOrderMenu(c *gin.Context) {
	// if shop not open redirect to /menu with an error
	isrunning, err := gorm.IsCafAppRunning()
	if err != nil || !isrunning {
		c.Redirect(http.StatusFound, "/menu")
	}

	data := make(map[string]interface{})
	data["Title"] = "Build Your Order"

	// check if user have any incomplete order
	user := getCurrentAuthUser(c)
	order, err := user.GetOneIncompleteOrder()
	if err != nil || order == nil || order.ID == 0 {
		logger.Info("Cannot get incomplete order from user. Assuming creating a fresh one.")
	} else {
		data["incompleteOrderURL"] = "/order/" + order.UUID
	}

	// get all menu items from db
	menu, err := gorm.GetActiveMenuItems()
	if err != nil {
		logger.Error("could not get menu items to display:", err)
		orderError(c, "Could not load menu items")
		return
	}
	data["menu"] = menu

	// render
	renderHTML(c, 200, "order-menu.html", data)
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
		return
	}

	// user's current balance
	data["balance"] = user.CurrentBalanceInCents

	// user's phone number
	data["phone"] = user.PhoneNumber

	// does user have gus id
	if user.GusID == 0 {
		data["needGusID"] = true
	}

	// does user have a phone number
	if user.PhoneNumber == "" {
		data["needPhoneNumber"] = true
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
			data["selectedDrinkName"] = order.OrderRows[i].Product.DisplayName
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
		return
	}
	data["menu"] = menu

	// get all drinks
	drinks, err := gorm.GetAllAddonProducts()
	if err != nil {
		logger.Error("could not get drinks to display:", err)
		orderError(c, "Database Error")
		return
	}
	data["drinks"] = drinks

	// get all destinations
	dest, err := gorm.GetAllDestinations()
	if err != nil {
		logger.Error("could not load destinations:", err)
		orderError(c, "Database Error")
		return
	}
	data["destinations"] = dest

	renderHTML(c, 200, "order-info.html", data)
}

func getFinalize(c *gin.Context, order gorm.Order) {
	data := make(map[string]interface{})
	data["Title"] = "Build Your Order"

	// determine currently selected meal and drink
	for i := range order.OrderRows {
		if order.OrderRows[i].RowType == gorm.RowTypeNormal {
			data["selectedMealName"] = order.OrderRows[i].Product.DisplayName
			data["selectedMealPrice"] = order.OrderRows[i].Product.PriceInCents
		}
		if order.OrderRows[i].RowType == gorm.RowTypeIncluded {
			data["selectedDrinkName"] = order.OrderRows[i].Product.DisplayName
		}
	}

	// pass in the order details
	data["deliveryFee"] = order.DeliveryFeeInCents
	data["orderTotal"] = order.TotalInCents
	data["cafAccountChargeAmount"] = order.CafAccountChargeAmountInCents

	// dest
	var dest gorm.Destination
	err := dest.PopulateByTag(order.DestinationTag)
	if err != nil {
		logger.Error(err)
		orderError(c, "Database error")
		return
	}
	data["destination"] = dest.Name

	// error if any
	sesh := sessions.Default(c)
	errText := getStringFromSession(sesh, "error")
	if errText != "" {
		sesh.Delete("error")
		sesh.Save()
		data["error"] = errText
	}

	renderHTML(c, 200, "order-finalize.html", data)
}

func getComplete(c *gin.Context, order gorm.Order) {
	data := make(map[string]interface{})
	data["Title"] = "Woohoo!"
	renderHTML(c, 200, "order-placed.html", data)
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
	// get from POST form
	selectedMeal := c.PostForm("meal")
	selectedDrink := c.PostForm("drink")
	selectedDestination := c.PostForm("destination")
	inputGusID := c.PostForm("gusID")
	inputPhoneNumber := c.PostForm("phone-input")
	// apply changes to meal
	if selectedMeal != "" {
		selectedMealInt, err := strconv.ParseUint(selectedMeal, 10, 32)
		if err != nil {
			logger.Error("invalid post form. Redirecting to edit page")
			orderError(c, "Bad Request. Bad. BAAADD")
			return
		}
		mealRow := order.GetMealRow() // potential nil pointer issue here but currently there's no way that happens
		var newMeal gorm.Product
		err = newMeal.PopulateByIDOnShelf(uint(selectedMealInt))
		if err != nil {
			logger.Error("cannot get meal product:", err)
			orderError(c, "Database error")
			return
		}
		mealRow.Product = &newMeal
		err = mealRow.Save()
		if err != nil {
			logger.Error("cannot save meal row:", err)
			orderError(c, "Database error")
			return
		}
	}

	// apply changes to drink
	if selectedDrink != "" {
		selectedDrinkInt, err := strconv.ParseUint(selectedDrink, 10, 32)
		if err != nil {
			logger.Error("invalid post form. Redirecting to edit page")
			orderError(c, "Bad Request. Bad. BAAADD")
			return
		}
		drinkRow := order.GetDrinkRow()
		var newDrink gorm.Product
		err = newDrink.PopulateByID(uint(selectedDrinkInt))
		if err != nil {
			logger.Error("cannot get drink product:", err)
			orderError(c, "Database error")
			return
		}
		if drinkRow != nil {
			drinkRow.Product = &newDrink
			err = drinkRow.Save()
			if err != nil {
				logger.Error("cannot save drink product:", err)
				orderError(c, "Database error")
				return
			}
		} else {
			newDrinkRow := gorm.NewOrderRowFromProduct(&newDrink)
			order.OrderRows = append(order.OrderRows, *newDrinkRow)
		}
	}

	// apply changes to destination
	if selectedDestination != "" {
		var dest gorm.Destination
		err := dest.PopulateByTag(selectedDestination)
		if err != nil {
			logger.Error("cannot get destination from db:", err)
			orderError(c, "Database error")
			return
		}
		order.DestinationTag = dest.Tag
	}

	// save user gus id if necessary
	user := getCurrentAuthUser(c)
	if user.GusID == 0 && inputGusID != "" {
		gusID, err := strconv.Atoi(inputGusID)
		if err != nil {
			logger.Error("invalid post form. Redirecting to edit page")
			orderError(c, "Bad Request. Bad. BAAADD")
			return
		}
		user.GusID = gusID
		err = user.Save()
		if err != nil {
			logger.Error("cannot save user gus id. Redirecting to edit page")
			orderError(c, "Bad Request. Bad. BAAADD")
			return
		}
		logger.Info("!!!!!! Gus User ID is", user.GusID)
	}

	//save user's phone
	if user.PhoneNumber == "" && inputPhoneNumber != "" {
		user.PhoneNumber = inputPhoneNumber
		err := user.Save()
		if err != nil {
			logger.Error("cannot save user phone. Redirecting to edit page")
			orderError(c, "Bad Request. Bad. BAAADD")
			return
		}
		logger.Info("!!!!!! Phone number is", user.GusID)
	}

	// save order
	order.CalculateDeliveryFee()
	order.CalculateTotal()
	order.StatusCode = gorm.OrderStatusFinalized
	err := order.Save()
	if err != nil {
		logger.Error("cannot save order:", err)
		orderError(c, "Database error")
		return
	}

	c.Redirect(http.StatusFound, "/order/"+order.UUID+"/finalize")
}

func postFinalize(c *gin.Context, order gorm.Order) {
	// is the user trying to go back to edit?
	tmp := c.PostForm("goToEdit")
	if tmp == "goToEdit" {
		order.StatusCode = gorm.OrderStatusIncomplete
		err := order.Save()
		if err != nil {
			logger.Error("cannot saving order")
			orderError(c, "Database error")
			return
		}
		c.Redirect(http.StatusFound, "/order/"+order.UUID)
		return
	}

	// check user's balance
	user := getCurrentAuthUser(c)
	if user.CurrentBalanceInCents < order.DeliveryFeeInCents {
		logger.Info("user have insufficient funds")
		sesh := sessions.Default(c)
		sesh.Set("error", "Error: You don't have enough funds in your CafApp balance for this delivery. You can add balance by <a class=\"text-info\" href=\"/redeem\">redeeming CafApp Delivery Cards</a>.")
		sesh.Save()
		c.Redirect(http.StatusFound, "/order/"+order.UUID)
		return
	}

	// subtract from user's balance
	user.CurrentBalanceInCents -= order.DeliveryFeeInCents
	err := user.Save()
	if err != nil {
		logger.Error("cannot update user's balance")
		orderError(c, "Fatal Error. Unable to place order. Your balance may have been wrongly subtracted. Please contact support for help.")
		return
	}

	// set status to Placed
	err = order.SetStatusTo(gorm.OrderStatusPlaced)
	if err != nil {
		logger.Error("cannot change status order")
		orderError(c, "Database error")
		return
	}

	// generate a tag
	err = order.GenerateTag()
	if err != nil {
		logger.Error("cannot generate tag for order")
		order.StatusCode = gorm.OrderStatusGeneralFailure
		order.Save()
		orderError(c, "Database error")
		return
	}

	// save
	err = order.Save()
	if err != nil {
		if err != nil {
			logger.Error("cannot save order")
			orderError(c, "Database error")
			return
		}
	}

	// signals the admin queue and order tracker
	socket.NewOrderPlaced(&order)

	c.Redirect(http.StatusFound, "/order/"+order.UUID+"/completed")
	return
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

// when frontend send an ajax request to redeem a delivery card
func handleRedeemDeliveryCard(c *gin.Context) {
	// bind
	type reqStruct struct {
		DeliveryCode string `json:"delivery_code"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// uppercase
	req.DeliveryCode = "CA-" + strings.ToUpper(req.DeliveryCode)

	// log
	logger.Info(req.DeliveryCode)

	// check if code exist in db
	code, err := gorm.GetRedeemableCodeByCode(req.DeliveryCode)
	if err == gorm.ErrRecordNotFound {
		logger.Error("code not found:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err != nil {
		logger.Error("unable to query db:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get user
	user := getCurrentAuthUser(c)

	// add the amount to user's balance
	ok, err := user.RedeemDeliveryCode(code)
	if err != nil {
		logger.Error("unable to redeem code:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !ok {
		logger.Error("invalid code entered:", req.DeliveryCode)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.String(200, "%v", user.CurrentBalanceInCents)
}
