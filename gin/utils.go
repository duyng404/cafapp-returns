package gin

import (
	"cafapp-returns/gorm"
	"cafapp-returns/logger"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// helper func to render html with all the necessary info for a template
func renderHTML(c *gin.Context, code int, template string, data map[string]interface{}) {
	// extracting common data
	u := getCurrentAuthUser(c)
	if u != nil {
		data["loggedIn"] = true
		data["currentUser"] = u
	}

	// write to a buffer
	buf, err := rdr.RenderHTML(template, data)
	if err != nil {
		logger.Error("error while generating html:", err)
		c.String(500, "500 Internal Server Error.")
		c.Abort()
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(code)
	c.Writer.Write(buf.Bytes())
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

// a helper func to use when error during code redeem.
// will redirect user to redeem page and display an err msg.
func redeemFailed(errmsg string, c *gin.Context) {
	session := sessions.Default(c)
	session.Set("error", errmsg)
	session.Save()
	c.Redirect(http.StatusFound, "/redeem")
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

func formatMoney(a int) string {
	l := a / 100
	r := a % 100
	ls := strconv.FormatUint(uint64(l), 10)
	rs := strconv.FormatUint(uint64(r), 10)
	if r < 10 {
		rs = "0" + rs
	}
	return "$" + ls + "." + rs
}

func rawHTML(s string) template.HTML {
	return template.HTML(s)
}

// helper func to get string from session
func getStringFromSession(s sessions.Session, name string) string {
	if a := s.Get(name); a != nil {
		if b, ok := a.(string); ok {
			return b
		}
	}
	return ""
}

// helper func to extract the last number of the order string.
// this will be referred to as the "friendly order number" and is the main way users will refer to the orders.
// example: SO-VSSP-43 will returns 43
func fromTagToNumber(tag string) string {
	splitted := strings.Split(tag, "-")
	return splitted[len(splitted)-1]
}

// get current user from gin store. If not logged in, will return a nil pointer
func getCurrentAuthUser(c *gin.Context) *gorm.User {
	if user, ok := c.Get("currentUser"); ok && user != nil {
		if user2, ok2 := user.(*gorm.User); ok2 && user2 != nil {
			return user2
		}
	}
	return nil
}

// calculate the delivery of a given order
// TODO: implement a proper rate
func calculateDeliveryFee(o *gorm.Order) int {
	return 75
}
