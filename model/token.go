package model

import (
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/random"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// TokenType is general token type (pwd recovery, verify email etc)
type TokenType int

// general token types
const (
	TokenTypePasswordRecovery TokenType = iota
	TokenTypeEmailVerification
)

func (tt TokenType) String() string {
	switch tt {
	case TokenTypePasswordRecovery:
		return "password_recovery"
	case TokenTypeEmailVerification:
		return "email_verification"
	default:
		return "unknown"
	}
}

const (
	tokenSize               = 64
	tokenMaxExpiryTimeHours = 24
)

var (
	msgInvalidToken = &i18n.Message{ID: "model.token.validate.app_error", Other: "invalid token"}
	msgTokenExpired = &i18n.Message{ID: "model.token.validate.expired.app_error", Other: "token has expired"}
)

// Token is an app token (verify email, pw recovery etc)
type Token struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

// NewToken returns the new token
func NewToken(tokentype TokenType, userID int64, expiryHours ...int) *Token {
	var expiry time.Duration

	if len(expiryHours) == 0 {
		expiry = time.Hour * tokenMaxExpiryTimeHours
	} else {
		expiry = time.Hour * time.Duration(expiryHours[0])
	}

	return &Token{
		UserID:    userID,
		Token:     random.SecureToken(tokenSize),
		Type:      tokentype.String(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expiry),
	}
}
