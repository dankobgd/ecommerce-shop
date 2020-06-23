package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

// AccessTokenLocation is the location from where the access token is retrieved
type AccessTokenLocation int

// access data information
const (
	AccessTokenType       = "bearer"
	AccessTokenGrantType  = "access_token"
	RefreshTokenGrantType = "refresh_token"

	AccessCookieName  = "access_token"
	RefreshCookieName = "refresh_token"

	HeaderBearer        = "Bearer"
	HeaderAuthorization = "Authorization"
)

// access token locations
const (
	TokenLocationNotFound AccessTokenLocation = iota
	TokenLocationCookie
	TokenLocationHeader
	TokenLocationQueryString
)

func (loc AccessTokenLocation) String() string {
	switch loc {
	case TokenLocationNotFound:
		return "not_found"
	case TokenLocationHeader:
		return "header"
	case TokenLocationCookie:
		return "cookie"
	case TokenLocationQueryString:
		return "query_string"
	default:
		return "unknown"
	}
}

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
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
	jwt.StandardClaims
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
