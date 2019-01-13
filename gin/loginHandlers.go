package gin

import (
	"cafapp-returns/ggoauth"
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Render the login button, when clicked will prompt user to log in through Google.
func handleGoogleLogin(c *gin.Context) {
	session := sessions.Default(c)
	// check if any pending error
	displayError := session.Get("error")
	if displayError != nil {
		session.Delete("error")
	}
	// generate a state and save it to the current session
	state := ggoauth.GenerateNewState()
	session.Set("state", state)
	session.Save()

	// render template, pass in the url to redirect the user after login
	c.HTML(200, "landing-gg-login.html", gin.H{
		"GGLoginUrl": ggoauth.GetLoginURL(state),
		"error":      displayError,
	})
}

// After logging in (or rejecting to log in), this endpoint will be called.
// Get user info from Google and tries to log them in to our system.
func handleGoogleLoginCallback(c *gin.Context) {
	// init sess, get state from sess
	session := sessions.Default(c)
	state := session.Get("state")
	session.Delete("state")
	session.Save()

	// validate state
	if state != c.Query("state") {
		logger.Error("Invalid session state")
		ggLoginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
		return
	}

	// get user info from google
	oauthResponse, err := ggoauth.GetUserDetailsFromGoogle(c.Query("code"))
	if err != nil {
		logger.Error("unable to get user info from google:", err)
		if err == ggoauth.ErrInvalidDomain {
			// not gustavus.edu? gtfo plz
			ggLoginFailed(fmt.Sprintf("You just logged in using the email %s. Please use your @gustavus.edu email!", oauthResponse.Email), c, session)
			return
		}
		// other errors
		ggLoginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
		return
	}

	// check if user already exist
	var user gorm.User
	err = user.PopulateByEmail(oauthResponse.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("error querying db for email", oauthResponse.Email, ":", err)
		ggLoginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
		return
	}
	if err == gorm.ErrRecordNotFound {
		// doesn't exist in our db, create new
		user.FullName = oauthResponse.Name
		user.FirstName = oauthResponse.GivenName
		user.LastName = oauthResponse.FamilyName
		user.Email = oauthResponse.Email
		splittedEmail := strings.Split(oauthResponse.Email, "@")
		user.GusUsername = splittedEmail[0]
		user.IsAdmin = false
		err = user.Create()
		if err != nil {
			logger.Error("error adding new user", user.Email, "to db:", err)
			ggLoginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
			return
		}
		// also save the oauth response
		googleUser := gorm.GoogleUser{
			OauthResponse: *oauthResponse,
			UserID:        user.ID,
		}
		err = googleUser.Create()
		if err != nil {
			logger.Error("error adding new *google* user", user.Email, "to db:", err)
			ggLoginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
			return
		}
	}

	// ok we got an user. generate a jwt for that user
	token, err := user.GenerateJWT()
	if err != nil {
		logger.Error("error generating jwt for user", user.Email, ":", err)
		ggLoginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
		return
	}

	// store the jwt in a cookie
	// cookie expire in 1 week
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
	})

	// login finished. redirect to next
	next := session.Get("next")
	if next == nil {
		next = "/"
	}
	// remember to unset next
	session.Delete("next")
	session.Save()
	c.Redirect(http.StatusFound, next.(string))

	return
}

// a helper func to use when error during login.
// will redirect user to login page and display an err msg.
func ggLoginFailed(errmsg string, c *gin.Context, session sessions.Session) {
	session.Set("error", errmsg)
	session.Save()
	c.Redirect(http.StatusFound, "/gg-login")
	return
}

// a helper func to use when accessing restricted pages without being logged in
// will save the current path in session, so after loggin in will be redirected
func stashThisPath(c *gin.Context, session sessions.Session) {
	// get current path, removing everything before the slash
	url := c.Request.URL
	url.Scheme = ""
	url.Opaque = ""
	url.User = nil
	url.Host = ""
	path := url.String()
	if path == "/gg-login" {
		return
	}
	session.Set("next", path)
	session.Save()
	c.Redirect(http.StatusFound, "/gg-login")
	return
}
