package gin

import (
	"cafapp-returns/config"
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAdminDash(c *gin.Context) {
	c.Redirect(http.StatusFound, config.AdminDashboardURL)
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
