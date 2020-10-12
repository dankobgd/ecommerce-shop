package model

import (
	"encoding/json"
	"io"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// error msgs
var (
	msgInvalidWishlist           = &i18n.Message{ID: "model.wishlist.validate.app_error", Other: "invalid wishlist data"}
	msgValidateWishlistUserID    = &i18n.Message{ID: "model.wishlist.validate.user_id.app_error", Other: "invalid wishlist user id"}
	msgValidateWishlistProductID = &i18n.Message{ID: "model.wishlist.validate.product_id.app_error", Other: "invalid wishlist product id"}
)

// Wishlist is the product wishlist
type Wishlist struct {
	ID        int64 `json:"id" db:"id"`
	UserID    int64 `json:"user_id" db:"user_id"`
	ProductID int64 `json:"product_id" db:"product_id"`
}

// WishlistFromJSON decodes the input and returns the Wishlist
func WishlistFromJSON(data io.Reader) (*Wishlist, error) {
	var p *Wishlist
	err := json.NewDecoder(data).Decode(&p)
	return p, err
}

// Validate validates the wishlist and returns an error if it doesn't pass criteria
func (pw *Wishlist) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if pw.UserID == 0 {
		errs.Add(Invalid("user_id", l, msgValidateWishlistUserID))
	}
	if pw.ProductID == 0 {
		errs.Add(Invalid("product_id", l, msgValidateWishlistProductID))
	}

	if !errs.IsZero() {
		return NewValidationError("Wishlist", msgInvalidWishlist, "", errs)
	}
	return nil
}
