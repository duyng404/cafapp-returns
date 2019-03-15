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
	// if shop not open redirect to /menu, but would allow if user is admin and AdminTestable is enabled
	user := getCurrentAuthUser(c)
	gvar, err := gorm.GetGlobalVar()
	if err != nil {
		c.Redirect(http.StatusFound, "/menu")
		return
	}
	isrunning, err := gorm.IsCafAppRunning()
	if user.IsAdmin && gvar.AdminTestable {
		isrunning = true
	}
	if err != nil || !isrunning {
		c.Redirect(http.StatusFound, "/menu")
		return
	}

	data := make(map[string]interface{})
	data["Title"] = "Build Your Order"

	// check if user have any incomplete order
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

	data["order"] = order

	// does user have gus id
	if order.User.GusID == 0 {
		data["needGusID"] = true
	}

	// does user have a phone number
	if order.User.PhoneNumber == "" {
		data["needPhoneNumber"] = true
	}

	// determine currently selected choices
	data["meal1id"] = 0
	data["side1id"] = 0
	data["drink1id"] = 0
	data["meal2id"] = 0
	data["side2id"] = 0
	data["drink2id"] = 0
	pseudototal := 0
	for i, v := range order.OrderRows {
		data["meal"+strconv.Itoa(i+1)+"id"] = v.MenuItemID
		pseudototal += v.MenuItem.DisplayPriceInCents
		for _, vv := range v.SubRows {
			if vv.Product.IsSide() {
				data["side"+strconv.Itoa(i+1)+"id"] = vv.ProductID
			}
			if vv.Product.IsDrink() {
				data["drink"+strconv.Itoa(i+1)+"id"] = vv.ProductID
			}
		}
	}

	// determine pseudo total
	pseudototal += order.DeliveryFeeInCents
	data["pseudototal"] = pseudototal
	data["cafAccountChargeAmount"] = pseudototal - order.DeliveryFeeInCents

	// get all menu items from db
	menu, err := gorm.GetActiveMenuItems()
	if err != nil {
		logger.Error("could not get menu items to display:", err)
		orderError(c, "Database Error")
		return
	}
	data["menu"] = menu

	// get all drinks
	drinks, err := gorm.GetAllDrinkProducts()
	if err != nil {
		logger.Error("could not get drinks to display:", err)
		orderError(c, "Database Error")
		return
	}
	data["drinks"] = drinks

	sides, err := gorm.GetAllSideProducts()
	if err != nil {
		logger.Error("could not get sides to display:", err)
		orderError(c, "Database Error")
		return
	}
	data["sides"] = sides

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

	data["order"] = order

	// dest
	var dest gorm.Destination
	err := dest.PopulateByTag(order.DestinationTag)
	if err != nil {
		logger.Error(err)
		orderError(c, "Database error")
		return
	}
	data["destination"] = dest.Name
	data["pickupspot"] = dest.PickUpLocation

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
		orderError(c, "Database error")
		return
	}

	// get selected
	tmp := c.PostForm("selected-item")
	selected, err := strconv.Atoi(tmp)
	if err != nil {
		logger.Error("error getting selected meal ("+tmp+") from POST form", err)
		orderError(c, "Internal error")
		return
	}

	// get the menu item
	var mi gorm.MenuItem
	err = mi.PopulateByID(uint(selected))
	if err != nil {
		logger.Error("error getting menu item from db", err)
		orderError(c, "Database error")
		return
	}

	// create a new order
	order, err := user.NewOrderFromMenuItem(&mi)
	if err != nil {
		logger.Error("error creating new order", err)
		orderError(c, "Database error")
		return
	}

	logger.Info("created order with uuid", order.UUID)

	c.Redirect(http.StatusFound, "/order/"+order.UUID)
}

// POST step 2: user finished filling out info
func postOrderInfo(c *gin.Context, order gorm.Order) {
	// get from POST form
	selectedMeal1 := c.PostForm("meal1")
	selectedDrink1 := c.PostForm("drink1")
	selectedMeal2 := c.PostForm("meal2")
	selectedDrink2 := c.PostForm("drink2")
	selectedDestination := c.PostForm("destination")
	inputGusID := c.PostForm("gusID")
	inputPhoneNumber := c.PostForm("phone-input")
	// apply changes to order row 1
	if selectedMeal1 != "" {
		selectedMeal1Int, err := strconv.ParseUint(selectedMeal1, 10, 32)
		if err != nil {
			logger.Error("invalid post form. Redirecting to edit page")
			orderError(c, "Bad Request. Bad. BAAADD")
			return
		}
		selectedDrink1Int, err := strconv.ParseUint(selectedDrink1, 10, 32)
		if err != nil {
			logger.Error("invalid post form. Redirecting to edit page")
			orderError(c, "Bad Request. Bad. BAAADD")
			return
		}
		var newMenuItem gorm.MenuItem
		err = newMenuItem.PopulateByID(uint(selectedMeal1Int))
		if err != nil {
			logger.Error("cannot get menu item:", err)
			orderError(c, "Database error")
			return
		}
		var newDrink gorm.Product
		err = newDrink.PopulateByID(uint(selectedDrink1Int))
		if err != nil {
			logger.Error("cannot get drink product:", err)
			orderError(c, "Database error")
			return
		}
		err = order.OrderRows[0].SetMainSubRowTo(newMenuItem.StartingMain)
		if err != nil {
			logger.Error("cannot set main:", err)
			orderError(c, "Database error")
			return
		}
		err = order.OrderRows[0].SetSideSubRowTo(newMenuItem.StartingSide)
		if err != nil {
			logger.Error("cannot set side:", err)
			orderError(c, "Database error")
			return
		}
		err = order.OrderRows[0].SetDrinkSubRowTo(&newDrink)
		if err != nil {
			logger.Error("cannot set drink:", err)
			orderError(c, "Database error")
			return
		}
		order.OrderRows[0].MenuItem = &newMenuItem
		err = order.OrderRows[0].Save()
		if err != nil {
			logger.Error("cannot save row:", err)
			orderError(c, "Database error")
			return
		}
	}

	// apply changes to order row 2
	if selectedMeal2 != "" {
		selectedMeal2Int, err := strconv.ParseUint(selectedMeal2, 10, 32)
		if err != nil {
			logger.Error("invalid post form. Redirecting to edit page")
			orderError(c, "Bad Request. Bad. BAAADD")
			return
		}
		selectedDrink2Int, err := strconv.ParseUint(selectedDrink2, 10, 32)
		if err != nil {
			logger.Error("invalid post form. Redirecting to edit page")
			orderError(c, "Bad Request. Bad. BAAADD")
			return
		}
		var newMenuItem gorm.MenuItem
		err = newMenuItem.PopulateByID(uint(selectedMeal2Int))
		if err != nil {
			logger.Error("cannot get menu item:", err)
			orderError(c, "Database error")
			return
		}
		var newDrink gorm.Product
		err = newDrink.PopulateByID(uint(selectedDrink2Int))
		if err != nil {
			logger.Error("cannot get drink product:", err)
			orderError(c, "Database error")
			return
		}
		if len(order.OrderRows) == 1 {
			newRow := gorm.OrderRow{}
			newRow.Create()
			order.OrderRows = append(order.OrderRows, newRow)
		}
		err = order.OrderRows[1].SetMainSubRowTo(newMenuItem.StartingMain)
		if err != nil {
			logger.Error("cannot set main:", err)
			orderError(c, "Database error")
			return
		}
		err = order.OrderRows[1].SetSideSubRowTo(newMenuItem.StartingSide)
		if err != nil {
			logger.Error("cannot set side:", err)
			orderError(c, "Database error")
			return
		}
		err = order.OrderRows[1].SetDrinkSubRowTo(&newDrink)
		if err != nil {
			logger.Error("cannot set drink:", err)
			orderError(c, "Database error")
			return
		}
		order.OrderRows[1].MenuItem = &newMenuItem
		err = order.OrderRows[1].Save()
		if err != nil {
			logger.Error("cannot save row:", err)
			orderError(c, "Database error")
			return
		}
	}

	// if remove meal
	if selectedMeal2 == "" && len(order.OrderRows) == 2 {
		order.OrderRows[1].Delete()
		order.OrderRows = order.OrderRows[:len(order.OrderRows)-1]
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
		Meal1ID  uint `json:"meal1id"`
		Drink1ID uint `json:"drink1id"`
		Meal2ID  uint `json:"meal2id"`
		Drink2ID uint `json:"drink2id"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// build a temporary order just to calculate the price
	meal1 := gorm.MenuItem{}
	drink1 := gorm.Product{}
	meal2 := gorm.MenuItem{}
	drink2 := gorm.Product{}
	if req.Meal1ID != 0 {
		err = meal1.PopulateByID(req.Meal1ID)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	if req.Drink1ID != 0 {
		err = drink1.PopulateByID(req.Drink1ID)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	if req.Meal2ID != 0 {
		err = meal2.PopulateByID(req.Meal2ID)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	if req.Drink2ID != 0 {
		err = drink2.PopulateByID(req.Drink2ID)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	order := gorm.Order{
		OrderRows: []gorm.OrderRow{
			gorm.OrderRow{
				SubRows: []gorm.SubRow{
					gorm.SubRow{
						Product: meal1.StartingMain,
					},
					gorm.SubRow{
						Product: meal1.StartingSide,
					},
					gorm.SubRow{
						Product: &drink1,
					},
				},
			},
			gorm.OrderRow{
				SubRows: []gorm.SubRow{
					gorm.SubRow{
						Product: meal2.StartingMain,
					},
					gorm.SubRow{
						Product: meal2.StartingSide,
					},
					gorm.SubRow{
						Product: &drink2,
					},
				},
			},
		},
	}
	order.CalculateDeliveryFee()
	order.CalculateTotal()

	// return
	type resStruct struct {
		Meal1Name        string `json:"meal1name"`
		Main1Name        string `json:"main1name"`
		Meal1Price       string `json:"meal1price"`
		Side1Name        string `json:"side1name"`
		Drink1Name       string `json:"drink1name"`
		Meal2Name        string `json:"meal2name"`
		Main2Name        string `json:"main2name"`
		Meal2Price       string `json:"meal2price"`
		Side2Name        string `json:"side2name"`
		Drink2Name       string `json:"drink2name"`
		DeliveryFee      string `json:"delivery_fee"`
		OrderTotal       string `json:"order_total"`
		CafAcctChargeAmt string `json:"caf_acct_charge_amt"`
	}
	res := resStruct{
		DeliveryFee: formatMoney(order.DeliveryFeeInCents),
		// OrderTotal:       formatMoney(order.TotalInCents),
		OrderTotal:       formatMoney(meal1.DisplayPriceInCents + meal2.DisplayPriceInCents + order.DeliveryFeeInCents),
		CafAcctChargeAmt: formatMoney(meal1.DisplayPriceInCents + meal2.DisplayPriceInCents),
	}
	if meal1.ID != 0 {
		res.Meal1Name = meal1.DisplayName
		res.Main1Name = meal1.StartingMain.DisplayName
		res.Meal1Price = formatMoney(meal1.DisplayPriceInCents)
		res.Side1Name = meal1.StartingSide.DisplayName
	}
	if drink1.ID != 0 {
		res.Drink1Name = drink1.DisplayName
	}
	if meal2.ID != 0 {
		res.Meal2Name = meal2.DisplayName
		res.Main2Name = meal2.StartingMain.DisplayName
		res.Meal2Price = formatMoney(meal2.DisplayPriceInCents)
		res.Side2Name = meal2.StartingSide.DisplayName
	}
	if drink2.ID != 0 {
		res.Drink2Name = drink2.DisplayName
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
