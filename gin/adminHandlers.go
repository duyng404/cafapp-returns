package gin

import (
	"cafapp-returns/config"
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func handleAdminDash(c *gin.Context) {
	c.Redirect(http.StatusFound, config.AdminDashboardURL)
}

func handleAdminDashDriver(c *gin.Context) {
	c.Redirect(http.StatusFound, config.AdminDashboardURL+"/driver")
}

func handleAdminInfo(c *gin.Context) {
	user := getCurrentAuthUser(c)
	token, err := user.GenerateSocketToken()
	if err != nil {
		logger.Error("error generating token for user", user.GusUsername, ":", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"admin_name":     user.FullName,
		"admin_username": user.GusUsername,
		"socket_token":   token,
	})
}

func handleAdminGetDestinations(c *gin.Context) {
	dest, err := gorm.GetAllDestinations()
	if err != nil {
		logger.Error("unable to get destinations:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, dest)
}

func handleAdminGetProducts(c *gin.Context) {
	prods, err := gorm.GetAllProducts()
	if err != nil {
		logger.Error("unable to get products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, prods)
}

func handleAdminViewQueue(c *gin.Context) {
	orders, err := gorm.GetOrdersForAdminViewQueue()
	if err != nil {
		logger.Error("unable to get orders for admin:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func handleAdminViewAllRedeemableCodes(c *gin.Context) {
	codes, err := gorm.GetAllRedeemableCodes()
	if err != nil {
		logger.Error("error getting all redeemable codes for admin dash:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, codes)
}

func handleAdminGenerateRedeemableCodes(c *gin.Context) {
	// TODO: factor out this struct
	type reqStruct struct {
		Amount int    `json:"amount"`
		Reason string `json:"reason"`
	}
	var req reqStruct

	// bind
	err := c.Bind(&req)
	if err != nil {
		logger.Error("error reading request:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	codes, err := gorm.GenerateRedeemableCodes(req.Amount, req.Reason)
	if err != nil {
		logger.Error("error generating 5 redeemable codes:")
		c.JSON(http.StatusInternalServerError, codes)
		return
	}
	c.JSON(http.StatusOK, codes)
}

func handleAdminCafAppOnOff(c *gin.Context) {
	gvar, err := gorm.GetGlobalVar()
	if err != nil {
		logger.Error("cannot get global var")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":       gvar.IsCafAppRunning,
		"announcement": gvar.FrontpageAnnouncement,
	})
}

func handlePostAdminCafAppOnOff(c *gin.Context) {
	type reqStruct struct {
		SetTo string `json:"set_to"`
	}
	var req reqStruct

	// bind
	err := c.Bind(&req)
	if err != nil {
		logger.Error("error reading request:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	gvar, err := gorm.GetGlobalVar()
	if err != nil {
		logger.Error("cannot get global var")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	logger.Info("setting cafapp running status to", req.SetTo)

	if req.SetTo == "on" {
		err := gvar.TurnCafAppOn()
		if err != nil {
			logger.Error("cannot set running mode", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":       gvar.IsCafAppRunning,
			"announcement": gvar.FrontpageAnnouncement,
		})
	}

	if req.SetTo == "off" {
		err := gvar.TurnCafAppOff()
		if err != nil {
			logger.Error("cannot get running mode", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":       gvar.IsCafAppRunning,
			"announcement": gvar.FrontpageAnnouncement,
		})
	}
}

func handleAdminSetAnnouncement(c *gin.Context) {
	type reqStruct struct {
		SetTo string `json:"set_to"`
	}
	var req reqStruct

	// bind
	err := c.Bind(&req)
	if err != nil {
		logger.Error("error reading request:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	gvar, err := gorm.GetGlobalVar()
	if err != nil {
		logger.Error("cannot get global var")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	logger.Info("setting cafapp announcement to", req.SetTo)

	err = gvar.SetFrontpageAnnouncement(req.SetTo)
	if err != nil {
		logger.Error("cannot set announcement", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":       gvar.IsCafAppRunning,
		"announcement": gvar.FrontpageAnnouncement,
	})
}

func handleAdminMenuStatus(c *gin.Context) {
	gvar, err := gorm.GetGlobalVar()
	if err != nil {
		logger.Error("cannot get global var")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	menus, err := gorm.GetAllMenus()
	if err != nil {
		logger.Error("cannot get all menus")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active_menu_id": gvar.ActiveMenuID,
		"menus":          menus,
	})
}

func handlePostAdminMenuStatus(c *gin.Context) {
	type reqStruct struct {
		SetTo uint `json:"set_to"`
	}
	var req reqStruct

	// bind
	err := c.Bind(&req)
	if err != nil {
		logger.Error("error reading request:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	gvar, err := gorm.GetGlobalVar()
	if err != nil {
		logger.Error("cannot get global var")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	setTo := uint(1)
	if gvar.ActiveMenuID == 1 {
		setTo = 2
	}
	gvar.ActiveMenuID = setTo
	err = gvar.Save()
	if err != nil {
		logger.Error("cannot set active menu", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	menus, err := gorm.GetAllMenus()
	if err != nil {
		logger.Error("cannot get all menus")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active_menu_id": gvar.ActiveMenuID,
		"menus":          menus,
	})
}

func handleAdminViewOrdersLast12Hours(c *gin.Context) {
	orders, err := gorm.GetAllOrdersLast12hours()
	if err != nil {
		logger.Error("cannot get orders for display")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, orders)
}

func handleAdminViewOrders(c *gin.Context) {
	rawDate := c.Query("date")

	if rawDate != "" {
		date, err := time.Parse(time.RFC3339, rawDate)
		if err != nil {
			logger.Error("cannot parse date:", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		date = date.UTC()

		orders, err := gorm.GetAllOrdersForDay(date)
		if err != nil {
			logger.Error("cannot get orders for display")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		logger.Info("got", len(orders), "from db.")

		c.JSON(http.StatusOK, orders)
		return

	} else {
		orders, err := gorm.GetAllOrdersLast12hours()
		if err != nil {
			logger.Error("cannot get orders for display")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		logger.Info("got", len(orders), "from db.")

		c.JSON(http.StatusOK, orders)
		return
	}
}
