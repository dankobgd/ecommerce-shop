package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

// general token types
const (
	TokenTypePasswordRecovery  = "password_recovery"
	TokenTypeEmailVerification = "email_verification"
)

// access data information
const (
	AccessTokenType       = "bearer"
	AccessTokenGrantType  = "access_token"
	RefreshTokenGrantType = "refresh_token"

	AccessCookieName  = "access_token"
	RefreshCookieName = "refresh_token"
)

// RefreshToken is the user refresh token
type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

// AccessData holds the auth access info
type AccessData struct {
	AccessUUID string
	UserID     int64
}

// TokenMetadata holds the tokens details
type TokenMetadata struct {
	TokenType      string
	AccessToken    string
	RefreshToken   string
	AccessUUID     string
	RefreshUUID    string
	AccessExpires  time.Time
	RefreshExpires time.Time
}

// Claims is the custom claims for the jwt
type Claims struct {
	Authorized bool     `json:"authorized"`
	Username   string   `json:"username,omitempty"`
	Roles      []string `json:"roles,omitempty"`
	*jwt.StandardClaims
}

// ToJSON converts the refresh token to json string
func (t *RefreshToken) ToJSON() string {
	b, _ := json.Marshal(t)
	return string(b)
}

// RefreshTokenFromJSON decodes the input and returns the RefreshToken
func RefreshTokenFromJSON(data io.Reader) (*RefreshToken, error) {
	var t *RefreshToken
	err := json.NewDecoder(data).Decode(&t)
	return t, err
}
