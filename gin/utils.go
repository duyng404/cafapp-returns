package gin

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// helper func to render html with all the necessary info for a template
func renderHTML(c *gin.Context, template string, data map[string]interface{}) {
	// extracting common data
	s := sessions.Default(c)
	if tmp := s.Get("loggedIn"); tmp != nil {
		if loggedIn, ok := tmp.(bool); ok && loggedIn {
			data["loggedIn"] = true
		}
	}
	c.HTML(200, template, data)
}

// a helper func to use when error during login.
// will redirect user to login page and display an err msg.
func loginFailed(errmsg string, c *gin.Context, session sessions.Session) {
	session.Set("error", errmsg)
	session.Save()
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
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
	return
}

// a helper func to use after user is logged in
// will redirect to next, if next is empty, redirect to homepage
func redirectToNext(c *gin.Context) {
	s := sessions.Default(c)

	next := s.Get("next")
	if next == nil {
		next = "/"
	}

	// remember to unset next
	s.Delete("next")
	s.Save()
	c.Redirect(http.StatusFound, next.(string))

	return
}
