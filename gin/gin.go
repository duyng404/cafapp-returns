package gin

import (
	"cafapp-returns/config"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

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
