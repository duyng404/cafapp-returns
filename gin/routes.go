package gin

import (
	"cafapp-returns/logger"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var count = 0

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// InitRoutes : Creates all of the routes for our application and returns a router
func InitRoutes() *gin.Engine {

	router := gin.New()

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// CORS babyyy
	router.Use(corsMiddleware())

	// Sessions middleware
	store := cookie.NewStore([]byte("secret")) // TODO: change secret & possible refactor
	router.Use(sessions.Sessions("mysession", store))

	// static
	router.Use(static.Serve("/static", static.LocalFile("./static", true)))

	router.GET("/ping", handlePing)
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "landing-top.html", gin.H{})
	})
	router.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/gg-login")
	})
	router.GET("/gg-login", handleGoogleLogin)
	router.GET("/gg-login-cb", handleGoogleLoginCallback)
	router.GET("/about", func(c *gin.Context) {
		logger.Info(c.Request.URL)
		c.HTML(200, "landing-about.html", gin.H{})
	})
	router.GET("/func", func(c *gin.Context) {
		c.HTML(200, "landing-func.html", gin.H{})
	})

	return router
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "pong", "increaseme": 0})
}
