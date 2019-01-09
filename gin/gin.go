package gin

import (
	"cafapp-returns/config"
	"cafapp-returns/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

// JSONError : Error object for JSON return to frontend
type JSONError struct {
	Error string `json:"error"`
}

// Run : starts a gin server
func Run() {
	port := config.Port
	router.Run(":" + port)
}

// ReturnRouter : returns a pointer to the engine
func ReturnRouter() *gin.Engine {
	return router
}

// SetTestMode : set gin mode as test
func SetTestMode() {
	gin.SetMode(gin.TestMode)
}

func init() {
	router = InitRoutes()
	initViews(router)
}

// nope : quickly nope out of a gin handle, with error logging, and sending the response
func nope(c *gin.Context, err interface{}, code int, msg string) {
	logger.Error(msg, err)
	c.JSON(code, gin.H{
		"err": fmt.Sprintf(msg+": %v", err),
	})
}
