package model

import (
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// error msgs
var (
	msgInvalidTag           = &i18n.Message{ID: "model.product_tag.validate.app_error", Other: "invalid tag data"}
	msgValidateTagID        = &i18n.Message{ID: "model.product_tag.validate.id.app_error", Other: "invalid  tag id"}
	msgValidateTagProductID = &i18n.Message{ID: "model.product_tag.validate.product_id.app_error", Other: "invalid tag product id"}
	msgValidateTagName      = &i18n.Message{ID: "model.product_tag.validate.name.app_error", Other: "invalid tag name"}
	msgValidateTagCrAt      = &i18n.Message{ID: "model.product_tag.validate.created_at.app_error", Other: "invalid tag created_at timestamp"}
	msgValidateTagUpAt      = &i18n.Message{ID: "model.product_tag.validate.updated_at.app_error", Other: "invalid tag updated_at timestamp"}
)

// ProductTag is the product tag
type ProductTag struct {
	ID        int64     `json:"id" db:"id" schema:"-"`
	ProductID int64     `json:"-" db:"product_id" schema:"-"`
	Name      string    `json:"name" db:"name" schema:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" schema:"-"`
}

// PreSave will fill timestamps
func (pt *ProductTag) PreSave() {
	pt.CreatedAt = time.Now()
	pt.UpdatedAt = pt.CreatedAt
}

// PreUpdate sets the update timestamp
func (pt *ProductTag) PreUpdate() {
	pt.UpdatedAt = time.Now()
}

// Validate validates the tag and returns an error if it doesn't pass criteria
func (pt *ProductTag) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if pt.ID != 0 {
		errs.Add(Invalid("tag.id", l, msgValidateTagID))
	}
	if pt.ProductID != 0 {
		errs.Add(Invalid("tag.product_id", l, msgValidateTagProductID))
	}
	if pt.Name == "" {
		errs.Add(Invalid("tag.name", l, msgValidateTagName))
	}
	if pt.CreatedAt.IsZero() {
		errs.Add(Invalid("tag.created_at", l, msgValidateTagCrAt))
	}
	if pt.UpdatedAt.IsZero() {
		errs.Add(Invalid("tag.updated_at", l, msgValidateTagUpAt))
	}

	if !errs.IsZero() {
		return NewValidationError("ProductTag", msgInvalidTag, "", errs)
	}
	return nil
}
