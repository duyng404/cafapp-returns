package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/jwt"
	"cafapp-returns/logger"
	"net/http"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

// authMiddleware : make sure restricted paths are protected
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok := checkJWT(c)
		if !ok {
			s := sessions.Default(c)
			logger.Info("JWT check failed. Assuming not logged in. Bouncing back to /login")
			stashThisPath(c, s)
			loginFailed("Please log in before accessing this page", c, s)
			return
		}

		c.Next()
	}
}

// if an user is logged in (indicated by the JWT in the cookie), then we check the JWT
// and extract the user for all handlers to use
func loginDetector() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok := checkJWT(c)
		if !ok {
			logger.Info("No JWT detected. A Guest!")
		} else {
			u := getCurrentAuthUser(c)
			logger.Info("Browsing as user", u.Email)
		}

		c.Next()
	}
}

// make sure restricted path are restricted and only accessible by admin users
func adminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check proper jwt
		ok := checkJWT(c)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// check if admin
		user := getCurrentAuthUser(c)
		if user.IsAdmin == false {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func checkJWT(c *gin.Context) bool {
	// get the raw jwt from cookie
	tokenString, err := c.Cookie("auth")
	if err != nil {
		logger.Info("error getting jwt from cookies:", err, ". Assuming not logged in.")
		return false
	}

	if tokenString == "" {
		logger.Warning("blank cookie. assuming not logged in")
		return false
	}

	// validate and stuff
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		logger.Warning("error parsing jwt:", err)
		return false
	}

	// we have the claims. does this user actually exist in db?
	user := gorm.User{}
	err = user.PopulateByID(claims.UserID)
	if err != nil {
		logger.Error("error looking up user from jwt:", err)
		return false
	}

	// save the user and the claim to gin's storage
	c.Set("currentUser", &user)
	c.Set("claims", claims)

	return true
}
