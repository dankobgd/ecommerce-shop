package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

// access data information
const (
	AccessTokenType       = "bearer"
	AccessTokenGrantType  = "access_token"
	RefreshTokenGrantType = "refresh_token"

	AccessCookieName  = "access_token"
	RefreshCookieName = "refresh_token"
)

// AccessToken is the user access token
type AccessToken struct {
	ID       string `json:"id"`
	Token    string `json:"token,omitempty"`
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
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

// UserClaims is the custom claims for the jwt
type UserClaims struct {
	Authorized bool     `json:"authorized"`
	Username   string   `json:"username,omitempty"`
	Roles      []string `json:"roles,omitempty"`
	*jwt.StandardClaims
}

// PreSave sets the uuid and isActive flag
func (t *AccessToken) PreSave() {
	t.ID = uuid.New().String()
	t.IsActive = true
}

// ToJSON converts token to json string
func (t *AccessToken) ToJSON() string {
	b, _ := json.Marshal(t)
	return string(b)
}

// AccessTokenFromJSON decodes the input and returns the AccessToken
func AccessTokenFromJSON(data io.Reader) (*AccessToken, error) {
	var t *AccessToken
	err := json.NewDecoder(data).Decode(&t)
	return t, err
}
