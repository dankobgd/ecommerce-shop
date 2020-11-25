package model

import (
	"encoding/json"
	"io"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// error msgs
var (
	msgInvalidProductTag           = &i18n.Message{ID: "model.product_tag.validate.app_error", Other: "invalid tag data"}
	msgValidateProductTagID        = &i18n.Message{ID: "model.product_tag.validate.id.app_error", Other: "invalid id"}
	msgValidateProductTagTagID     = &i18n.Message{ID: "model.product_tag.validate.tag_id.app_error", Other: "invalid tag_id"}
	msgValidateProductTagProductID = &i18n.Message{ID: "model.product_tag.validate.product_id.app_error", Other: "invalid product_id"}
)

// ProductTag is the product tag
type ProductTag struct {
	TagID     *int64 `json:"tag_id" db:"tag_id"`
	ProductID *int64 `json:"product_id" db:"product_id"`
	*Tag
}

// ProductTagPatch is the patch for product tag
type ProductTagPatch struct {
	TagID *int64 `json:"tag_id"`
}

// ProductTagFromJSON decodes the input and returns the ProductTag
func ProductTagFromJSON(data io.Reader) (*ProductTag, error) {
	var pt *ProductTag
	err := json.NewDecoder(data).Decode(&pt)
	return pt, err
}

// Patch patches the product tag
func (pt *ProductTag) Patch(patch *ProductTagPatch) {
	if patch.TagID != nil {
		pt.TagID = patch.TagID
	}
}

// ProductTagPatchFromJSON decodes the input and returns the ProductTagPatch
func ProductTagPatchFromJSON(data io.Reader) (*ProductTagPatch, error) {
	var p *ProductTagPatch
	err := json.NewDecoder(data).Decode(&p)
	return p, err
}

// Validate validates the tag and returns an error if it doesn't pass criteria
func (pt *ProductTag) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if pt.TagID != nil {
		errs.Add(Invalid("tag_id", l, msgValidateProductTagTagID))
	}
	if pt.ProductID != nil {
		errs.Add(Invalid("product_id", l, msgValidateProductTagProductID))
	}

	if !errs.IsZero() {
		return NewValidationError("ProductTag", msgInvalidProductTag, "", errs)
	}
	return nil
}
