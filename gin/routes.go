package gin

import (
	"cafapp-returns/ggoauth"
	"cafapp-returns/logger"
	"fmt"
	"io/ioutil"
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
	router.GET("/gg-login", func(c *gin.Context) {
		state := ggoauth.GenerateNewState()
		session := sessions.Default(c)
		session.Set("state", state)
		session.Save()
		c.HTML(200, "landing-gg-login.html", gin.H{
			"GGLoginUrl": ggoauth.GetLoginURL(state),
		})
	})
	router.GET("/gg-login-cb", func(c *gin.Context) {
		session := sessions.Default(c)
		state := session.Get("state")
		if state != c.Query("state") {
			logger.Error("Invalid session state")
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state"))
			return
		}
		token, err := ggoauth.GetTokenFromCode(c.Query("code"))
		if err != nil {
			logger.Error(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		client := ggoauth.GetClientFromToken(token)
		email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			logger.Error(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer email.Body.Close()
		data, _ := ioutil.ReadAll(email.Body)
		logger.Info("email body: ", string(data))
		c.Status(200)
	})
	router.GET("/about", func(c *gin.Context) {
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
