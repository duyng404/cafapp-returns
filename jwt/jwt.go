package jwt

import (
	"cafapp-returns/config"
	"cafapp-returns/logger"
	"errors"

	jwtgo "github.com/dgrijalva/jwt-go"
)

var (
	// ErrInvalidToken when the token is invalid
	ErrInvalidToken = errors.New("Invalid Token")
)

// Claims are the structure of the jwt
type Claims struct {
	UserID      uint
	GusUsername string
	IsAdmin     bool
	// other necessary jwt claim info
	jwtgo.StandardClaims
}

// NewToken creates a new jwt token for the given parameters
func NewToken(userID uint, gusUsername string, isAdmin bool, expiration int64) (string, error) {
	claims := Claims{
		UserID:      userID,
		GusUsername: gusUsername,
		IsAdmin:     isAdmin,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: expiration,
			Issuer:    "cafapp",
		},
	}
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSigningKey))
}

// ParseToken parses and check a jwt token if its valid
func ParseToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwtgo.ParseWithClaims(tokenString, &Claims{}, func(token *jwtgo.Token) (interface{}, error) {
		return config.JWTSigningKey, nil
	})
	if err != nil {
		logger.Error("error while parsing jwt token:", err)
		return nil, ErrInvalidToken
	}

	// validate
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// probably wont be, but still, check if userid is zero
	if token.Claims.(*Claims).UserID == 0 {
		logger.Error("id in jwt token is zero wtf")
		return nil, ErrInvalidToken
	}

	return token.Claims.(*Claims), err
}
