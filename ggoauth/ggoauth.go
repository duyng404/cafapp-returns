package ggoauth

import (
	"cafapp-returns/config"
	"cafapp-returns/logger"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http/httputil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	conf *oauth2.Config
	// ErrInvalidCode when the code is invalid
	ErrInvalidCode = errors.New("invalid code, unable to get token")
	// ErrInvalidDomain when it's not gustavus.edu
	ErrInvalidDomain = errors.New("invalid domain")
)

// OauthResponse the response containing user info from Google, used to identify user
type OauthResponse struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email" gorm:"index:email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
	HostDomain    string `json:"hd"`
}

func init() {
	conf = &oauth2.Config{
		ClientID:     config.GGoauthID,
		ClientSecret: config.GGoauthSecret,
		RedirectURL:  config.GGoauthRedirectURL,
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

// GenerateNewState simply generate a base64 random string, to be used as state
func GenerateNewState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// GetLoginURL from a given state generate an url that is to be attached to the google login button
func GetLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

// GetUserDetailsFromGoogle takes in the code provided by google after user
// logged in, verify it and get the user info from google
func GetUserDetailsFromGoogle(code string) (*OauthResponse, error) {
	// get token from code
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		logger.Error(err)
		return nil, ErrInvalidCode
	}

	// create a http client from token
	client := conf.Client(oauth2.NoContext, token)
	// make a GET request to get user info
	rawUser, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return nil, err
	}
	// remember to close
	defer rawUser.Body.Close()

	// log the info
	t, _ := httputil.DumpResponse(rawUser, true)
	logger.Info("Got an user from Google: ")
	logger.Info(string(t))

	oauthResponse := OauthResponse{}
	err = json.NewDecoder(rawUser.Body).Decode(&oauthResponse)
	if err != nil {
		return nil, err
	}

	if oauthResponse.HostDomain != "gustavus.edu" {
		return &oauthResponse, ErrInvalidDomain
	}

	return &oauthResponse, nil
}
