package ggoauth

import (
	"cafapp-returns/config"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf *oauth2.Config

func init() {
	conf = &oauth2.Config{
		ClientID:     config.GGoauthID,
		ClientSecret: config.GGoauthSecret,
		RedirectURL:  "http://localhost:7000/gg-login-cb",
		Scopes:       []string{"profile", "email", "openid"},
		Endpoint:     google.Endpoint,
	}
}

func GenerateNewState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func GetLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func GetTokenFromCode(code string) (*oauth2.Token, error) {
	return conf.Exchange(oauth2.NoContext, code)
}

func GetClientFromToken(token *oauth2.Token) *http.Client {
	return conf.Client(oauth2.NoContext, token)
}
