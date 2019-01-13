package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/jwt"
	"cafapp-returns/logger"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware :
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// get the raw jwt from cookie
		tokenString, err := c.Cookie("auth")
		if err != nil {
			logger.Error("error getting jwt from cookies:", err)
			stashThisPath(c, session)
			ggLoginFailed("Please log in before accessing this page", c, session)
			return
		}

		if tokenString == "" {
			logger.Error("blank cookie")
			stashThisPath(c, session)
			ggLoginFailed("Please log in before accessing this page", c, session)
			return
		}

		// validate and stuff
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			logger.Error("error parsing jwt:", err)
			stashThisPath(c, session)
			ggLoginFailed("Please log in again to access", c, session)
			return
		}

		// we have the claims. does this user actually exist in db?
		user := gorm.User{}
		err = user.PopulateById(claims.UserID)
		if err != nil {
			logger.Error("error looking up user from jwt:", err)
			stashThisPath(c, session)
			ggLoginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
			return
		}

		// save the user and the claim to gin's storage
		c.Set("currentUser", &user)
		c.Set("claims", claims)

		c.Next()
	}
}
