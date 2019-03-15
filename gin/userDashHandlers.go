package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func handleUserDash(c *gin.Context) {
	data := make(map[string]interface{})
	data["Title"] = "My Account"

	//display all past orders
	user := getCurrentAuthUser(c)
	orders, err := gorm.GetAllOrderFromUser(user.ID)
	if err != nil {
		logger.Error("Cannot display past orders", err)
		return
	}
	data["orders"] = orders

	//current user info
	data["user"] = user
	data["totalOrders"] = len(*orders)
	renderHTML(c, 200, "userdash-top.html", data)
}

func handleOrderDetail(c *gin.Context) {
	data := make(map[string]interface{})
	data["Title"] = "Order Details"
	var order gorm.Order
	stuff := c.Param("stuff")
	err := order.PopulateByUUID(stuff)
	if err != nil {
		logger.Error("order uuid in params is not valid:", err)
		logger.Info("boucing back to /dash")
		c.Redirect(http.StatusFound, "/dash")
	}

	if order.StatusCode == gorm.OrderStatusDelivered {
		data["isDelivered"] = true
		deliver, err := gorm.GetDeliveredTime(order.ID)
		if err != nil {
			logger.Error("cannot display delivery time of order")
			return
		}
		data["time"] = deliver
	}

	for i := range order.OrderRows {
		if order.OrderRows[i].RowType == gorm.RowTypeNormal {
			data["selectedMealName"] = order.OrderRows[i].Product.DisplayName
			data["selectedMealPrice"] = order.OrderRows[i].Product.PriceInCents
		}
		if order.OrderRows[i].RowType == gorm.RowTypeIncluded {
			data["selectedDrinkName"] = order.OrderRows[i].Product.DisplayName
		}
	}
	var dest gorm.Destination
	err2 := dest.PopulateByTag(order.DestinationTag)
	if err2 != nil {
		logger.Error(err2)
		orderError(c, "Database error")
		return
	}

	data["destination"] = dest.Name
	data["deliveryFee"] = order.DeliveryFeeInCents
	data["orderTotal"] = order.TotalInCents
	data["cafAccountChargeAmount"] = order.CafAccountChargeAmountInCents
	renderHTML(c, 200, "userdash-detail.html", data)
}

func handleUserRedeem(c *gin.Context) {
	data := make(map[string]interface{})
	data["Title"] = "Redeem"

	// get user
	user := getCurrentAuthUser(c)
	data["User"] = user

	// get pending errors if any
	session := sessions.Default(c)
	err := session.Get("error")
	session.Delete("error")
	session.Save()
	data["Error"], _ = err.(string)

	renderHTML(c, 200, "userdash-redeem.html", data)
}

func handleUserRedeemPost(c *gin.Context) {
	// get the code from input
	input := c.PostForm("redeem-input")
	input = strings.ToUpper(input)
	input = "CA-" + input

	// check if code exist in db
	code, err := gorm.GetRedeemableCodeByCode(input)
	if err == gorm.ErrRecordNotFound {
		logger.Error("code not found:", err)
		redeemFailed("You have entered an invalid code.", c)
		return
	}
	if err != nil {
		logger.Error("unable to query db:", err)
		redeemFailed("Oh no! Some unexpected error happened. Please try again later. Your code was not redeemed, and should still be available.", c)
		return
	}

	// get current logged in user
	user := getCurrentAuthUser(c)

	// add the amount to user's balance
	ok, err := user.RedeemDeliveryCode(code)
	if err != nil {
		logger.Error("unable to redeem code:", err)
		redeemFailed("Oh no! Some technical error happened. Please contact our support teams for assistance.", c)
		return
	}
	if !ok {
		logger.Error("invalid code entered:", input)
		redeemFailed("You have entered an invalid code.", c)
		return
	}

	session := sessions.Default(c)
	session.Set("redeemSuccess", true)
	session.Set("newBalance", user.CurrentBalanceInCents)
	session.Save()
	c.Redirect(http.StatusFound, "/redeem-success")
}

func handleUserRedeemSuccess(c *gin.Context) {
	data := make(map[string]interface{})
	data["Title"] = "Woohoo!"

	session := sessions.Default(c)
	success, ok := session.Get("redeemSuccess").(bool)
	if !ok || !success {
		c.Redirect(http.StatusFound, "/redeem")
		return
	}
	newBalance, ok := session.Get("newBalance").(int)
	if !ok {
		user := getCurrentAuthUser(c)
		newBalance = user.CurrentBalanceInCents
	}
	data["NewBalance"] = newBalance

	session.Delete("redeemSuccess")
	session.Delete("newBalance")
	session.Save()

	renderHTML(c, 200, "userdash-redeem-success.html", data)
}

func handleEditPhoneNumbers(c *gin.Context) {
	// bind
	type reqStruct struct {
		Phone string `json:"phone"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := getCurrentAuthUser(c)

	//save to db
	err = user.SaveUserPhone(req.Phone, user.ID)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(200, user.PhoneNumber)
}
