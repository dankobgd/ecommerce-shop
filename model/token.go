package model

import (
	"net/http"
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/utils/random"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

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
func NewToken(tokentype string, userID int64) *Token {
	return &Token{
		UserID:    userID,
		Token:     random.SecureToken(tokenSize),
		Type:      tokentype,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * tokenMaxExpiryTimeHours),
	}
}

func newInvalidTokenError(msg *i18n.Message, errs ValidationErrors) *AppErr {
	details := map[string]interface{}{}
	if !errs.IsZero() {
		details["validation"] = map[string]interface{}{"errors": errs}
	}
	return NewAppErr("Token.Validate", ErrInvalid, locale.GetUserLocalizer("en"), msg, http.StatusUnprocessableEntity, details)
}

// Validate validates the token
func (t *Token) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if len(t.Token) != tokenSize {
		errs.Add(NewValidationErr("token", l, msgInvalidToken))
	}
	if t.CreatedAt.IsZero() {
		errs.Add(NewValidationErr("token", l, msgTokenExpired))
	}

	if !errs.IsZero() {
		return newInvalidTokenError(msgInvalidToken, errs)
	}
	return nil
}
