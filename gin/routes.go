package gin

import (
	"cafapp-returns/config"
	"cafapp-returns/socket"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var count = 0

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.AdminDashboardURL)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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
	store := cookie.NewStore([]byte(config.SessionCookieKey))
	store.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   604800, // a week
		Path:     "/",
	})
	router.Use(sessions.Sessions("mysession", store))

	// login detector
	router.Use(loginDetector())

	// static
	router.Use(static.Serve("/static", static.LocalFile("./static", true)))
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	// ping
	router.GET("/ping", handlePing)

	//404
	router.NoRoute(func(c *gin.Context) {
		renderHTML(c, 404, "landing-404.html", gin.H{})
	})

	// landing group contains public-facing paths, aka, anyone can see without logging in
	landing := router.Group("/")
	{
		landing.GET("/", handleLandingTop)
		landing.GET("/about", handleLandingAbout)
		landing.GET("/news", handleLandingNews)
		landing.GET("/menu", handleLandingMenu)
		landing.GET("/faq", handleLandingFAQ)
	}

	// login group will handle logging users in and out
	login := router.Group("/")
	{
		login.GET("/login", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/gg-login")
		})
		login.GET("/gg-login", handleGoogleLogin)
		login.GET("/gg-login-cb", handleGoogleLoginCallback)
		login.GET("/logout", handleLogout)
	}

	// restricted group requires loggin in before accessing
	restricted := router.Group("/", authMiddleware())
	{
		restricted.GET("/order", handleOrderGet)
		restricted.GET("/order/:stuff", handleOrderGet)
		restricted.GET("/order/:stuff/:action", handleOrderGet)
		restricted.POST("/order", handleOrderPost)
		restricted.POST("/order/:stuff", handleOrderPost)
		restricted.POST("/order/:stuff/:action", handleOrderPost)

		restricted.GET("/dash", handleUserDash)
		restricted.GET("/dash/order/:orderuuid", handleOrderDetail)
		restricted.GET("/redeem", handleUserRedeem)
		restricted.GET("/redeem-success", handleUserRedeemSuccess)
		restricted.POST("/redeem", handleUserRedeemPost)

		restricted.GET("/tracker", handleOrderTracker)
		restricted.GET("/admin", handleAdminDash)
		restricted.GET("/driver", handleAdminDashDriver)
	}

	// api group for frontend interaction, will require auth
	api := router.Group("/api", authMiddleware())
	{
		api.POST("/recalculate-order", handleRecalculateOrder)
		api.GET("/my-info", handleUserInfo)
		api.POST("/quick-redeem", handleRedeemDeliveryCard)
		api.POST("/edit-phone", handleEditPhoneNumbers)
		api.GET("/view-active-orders", handleViewActiveOrders)
	}

	// api group for admin dash, will require auth with admin privilege
	apiadmin := router.Group("/api/admin", adminAuthMiddleware())
	{
		apiadmin.GET("/my-info", handleAdminInfo)
		apiadmin.GET("/view-queue", handleAdminViewQueue)
		apiadmin.GET("/destination", handleAdminGetDestinations)
		apiadmin.GET("/product", handleAdminGetProducts)
		apiadmin.GET("/view-redeemable-codes", handleAdminViewAllRedeemableCodes)
		apiadmin.GET("/view-users", handleAdminViewUsers)
		apiadmin.GET("/view-users/:userid", handleAdminViewOneUser)
		apiadmin.POST("/generate-redeemable-codes", handleAdminGenerateRedeemableCodes)
		apiadmin.GET("/cafapp-onoff", handleAdminCafAppOnOff)
		apiadmin.POST("/cafapp-onoff", handlePostAdminCafAppOnOff)
		apiadmin.GET("/menu-status", handleAdminMenuStatus)
		apiadmin.POST("/menu-status", handlePostAdminMenuStatus)
		apiadmin.GET("/orders-last-12-hours", handleAdminViewOrdersLast12Hours)
		apiadmin.GET("/view-orders", handleAdminViewOrders)
	}

	// TODO: make a group for this and look at authentication
	router.GET("/socket/", gin.WrapH(socket.GetServer()))
	router.POST("/socket/", gin.WrapH(socket.GetServer()))
	return router
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "pong", "increaseme": 0})
}
