package model

import (
	"encoding/json"
	"io"
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
	ID        *int64     `json:"id" db:"tag_id" schema:"-"`
	ProductID *int64     `json:"product_id" db:"tag_product_id" schema:"-"`
	Name      *string    `json:"name" db:"tag_name" schema:"name"`
	CreatedAt *time.Time `json:"created_at" db:"tag_created_at" schema:"-"`
	UpdatedAt *time.Time `json:"updated_at" db:"tag_updated_at" schema:"-"`
}

// ProductTagPatch is the patch for product tag
type ProductTagPatch struct {
	Name *string `json:"name"`
}

// Patch patches the product tag
func (tag *ProductTag) Patch(patch *ProductTagPatch) {
	if tag.Name != nil {
		tag.Name = patch.Name
	}
}

// ProductTagPatchFromJSON decodes the input and returns the ProductTagPatch
func ProductTagPatchFromJSON(data io.Reader) (*ProductTagPatch, error) {
	var p *ProductTagPatch
	err := json.NewDecoder(data).Decode(&p)
	return p, err
}

// PreSave will fill timestamps
func (tag *ProductTag) PreSave() {
	now := time.Now()
	tag.CreatedAt = &now
	tag.UpdatedAt = tag.CreatedAt
}

// PreUpdate sets the update timestamp
func (tag *ProductTag) PreUpdate() {
	now := time.Now()
	tag.UpdatedAt = &now
}

// Validate validates the tag and returns an error if it doesn't pass criteria
func (tag *ProductTag) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if tag.ID != nil {
		errs.Add(Invalid("tag.id", l, msgValidateTagID))
	}
	if tag.ProductID != nil {
		errs.Add(Invalid("tag.product_id", l, msgValidateTagProductID))
	}
	if tag.Name == nil {
		errs.Add(Invalid("tag.name", l, msgValidateTagName))
	}
	if tag.CreatedAt.IsZero() {
		errs.Add(Invalid("tag.created_at", l, msgValidateTagCrAt))
	}
	if tag.UpdatedAt.IsZero() {
		errs.Add(Invalid("tag.updated_at", l, msgValidateTagUpAt))
	}

	if !errs.IsZero() {
		return NewValidationError("ProductTag", msgInvalidTag, "", errs)
	}
	return nil
}
