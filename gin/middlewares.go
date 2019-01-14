package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/jwt"
	"cafapp-returns/logger"

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
