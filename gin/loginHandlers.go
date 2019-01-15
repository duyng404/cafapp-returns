package gin

import (
	"cafapp-returns/ggoauth"
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Render the login button, when clicked will prompt user to log in through Google.
func handleGoogleLogin(c *gin.Context) {
	session := sessions.Default(c)
	// check if user already logged in
	ok := checkJWT(c)
	if ok {
		// already logged in, redirect to next
		redirectToNext(c)
		return
	}

	// check if any pending error
	var displayError string
	tmp := session.Get("error")
	if tmp != nil {
		session.Delete("error")
		if displayError, ok = tmp.(string); !ok {
			displayError = ""
		}
	}
	// generate a state and save it to the current session
	state := ggoauth.GenerateNewState()
	session.Set("state", state)
	session.Save()

	// render template, pass in the url to redirect the user after login
	renderHTML(c, 200, "landing-gg-login.html", gin.H{
		"GGLoginUrl": ggoauth.GetLoginURL(state),
		"error":      template.HTML(displayError),
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
		loginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
		return
	}

	// get user info from google
	oauthResponse, err := ggoauth.GetUserDetailsFromGoogle(c.Query("code"))
	if err != nil {
		logger.Error("unable to get user info from google:", err)
		if err == ggoauth.ErrInvalidDomain {
			// not gustavus.edu? gtfo plz
			loginFailed(fmt.Sprintf(`You attempted to log in as <strong>`+oauthResponse.Email+`</strong>. Please try again using your @gustavus.edu email.`), c, session)
			return
		}
		// other errors
		loginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
		return
	}

	// check if user already exist
	var user gorm.User
	err = user.PopulateByEmail(oauthResponse.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("error querying db for email", oauthResponse.Email, ":", err)
		loginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
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
			loginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
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
			logger.Info("gg user dump:", spew.Sdump(googleUser))
			// loginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
			// return
		}
	}

	// ok we got an user. generate a jwt for that user
	token, err := user.GenerateJWT()
	if err != nil {
		logger.Error("error generating jwt for user", user.Email, ":", err)
		loginFailed("Oh no! Login was unsuccessful. Maybe try again?", c, session)
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

	logger.Info(fmt.Sprintf("user %s just logged in", user.Email))

	// login finished. redirect to next
	redirectToNext(c)

	return
}

// logs the user out
func handleLogout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()

	// log out the user by clearing the auth cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	// all these for a simple line of logging "user xyz just logged out"
	ok := checkJWT(c)
	if !ok {
		logger.Error("weird error: no valid jwt when logging out an user")
	}
	if user, ok := c.Get("currentUser"); ok && user != nil {
		if user2, ok2 := user.(*gorm.User); ok2 && user2 != nil {
			logger.Info(fmt.Sprintf("user %s just logged out", user2.Email))
		}
	}

	c.Redirect(http.StatusFound, "/")
}
