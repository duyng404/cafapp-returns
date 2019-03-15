package gin

import (
	"cafapp-returns/config"
	"cafapp-returns/renderer"
	"html/template"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
	rdr    *renderer.Rdr
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
	// init render
	views := []string{
		"landing",
		"userdash",
		"order",
		"404",
		"tracker",
	}
	f := template.FuncMap{
		"formatMoney":      formatMoney,
		"rawHTML":          rawHTML,
		"fromTagToNumber":  fromTagToNumber,
		"statusCodeToText": statusCodeToText,
		"addOne":           addOne,
	}
	rdr = renderer.InitRdr(views, f)
	router = InitRoutes()
}
