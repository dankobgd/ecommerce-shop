package model

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgInvalidProduct      = &i18n.Message{ID: "model.product.validate.app_error", Other: "invalid product data"}
	msgValidateProductCrAt = &i18n.Message{ID: "model.product.validate.created_at.app_error", Other: "invalid created_at timestamp"}
	msgValidateProductUpAt = &i18n.Message{ID: "model.product.validate.updated_at.app_error", Other: "invalid updated_at timestamp"}
)

// Product represents the shop product model
type Product struct {
	ID          int64      `json:"id" db:"id"`
	BrandID     int        `json:"brand_id" db:"brand_id"`
	DiscountID  *int       `json:"discount_id" db:"discount_id"`
	Name        string     `json:"name" db:"name"`
	Slug        string     `json:"slug" db:"slug"`
	ImageURL    string     `json:"image_url" db:"image_url"`
	Description string     `json:"description" db:"description"`
	Price       int        `json:"price" db:"price"`
	Stock       int        `json:"stock" db:"stock"`
	SKU         string     `json:"sku" db:"sku"`
	IsFeatured  bool       `json:"is_featured" db:"is_featured"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}

// ProductFromJSON decodes the input and return the Product
func ProductFromJSON(data io.Reader) (*Product, error) {
	var p *Product
	err := json.NewDecoder(data).Decode(&p)
	return p, err
}

// ToJSON converts Product to json string
func (p *Product) ToJSON() string {
	b, _ := json.Marshal(p)
	return string(b)
}

// PreSave will set missing defaults and fill CreatedAt and UpdatedAt times
func (p *Product) PreSave() {
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
}

// PreUpdate sets the update timestamp
func (p *Product) PreUpdate() {
	p.UpdatedAt = time.Now()
}

// NewInvalidProductError builds the invalid product error
func NewInvalidProductError(msg *i18n.Message, userID string, errs ValidationErrors) *AppErr {
	details := map[string]interface{}{}
	if !errs.IsZero() {
		details["validation"] = map[string]interface{}{"errors": errs}
	}
	return NewAppErr("Product.Validate", ErrInvalid, locale.GetUserLocalizer("en"), msg, http.StatusUnprocessableEntity, details)
}

// Validate validates the user and returns an error if it doesn't pass criteria
func (p *Product) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if p.ID != 0 {
		errs.Add(NewValidationErr("id", l, MsgValidateUserID))
	}
	if p.CreatedAt.IsZero() {
		errs.Add(NewValidationErr("created_at", l, msgValidateProductCrAt))
	}
	if p.UpdatedAt.IsZero() {
		errs.Add(NewValidationErr("updated_at", l, msgValidateProductUpAt))
	}

	if !errs.IsZero() {
		return NewInvalidUserError(msgInvalidProduct, "", errs)
	}
	return nil
}
