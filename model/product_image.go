package model

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"time"

	"github.com/dankobgd/ecommerce-shop/gocloudinary"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgValidateProductImage     = &i18n.Message{ID: "model.product_image.validate.url.app_error", Other: "invalid product image"}
	msgValidateProductImageSize = &i18n.Message{ID: "model.product_image.validate.url.app_error", Other: "File size exceeded, max 3MB allowed"}
)

// ProductImage is the product image
type ProductImage struct {
	ID        *int64     `json:"id" db:"id" schema:"-"`
	ProductID *int64     `json:"product_id" db:"product_id" schema:"-"`
	URL       *string    `json:"url" db:"url" schema:"-"`
	PublicID  *string    `json:"public_id" db:"public_id" schema:"-"`
	CreatedAt *time.Time `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at" schema:"-"`
}

// ProductImagePatch is the patch for tag
type ProductImagePatch struct {
	URL      *string `json:"url" schema:"-"`
	PublicID *string `json:"public_id" schema:"-"`
}

// ProductImagePatchFromJSON decodes the input and returns the ProductImagePatch
func ProductImagePatchFromJSON(data io.Reader) (*ProductImagePatch, error) {
	var p *ProductImagePatch
	err := json.NewDecoder(data).Decode(&p)
	return p, err
}

// Validate validates the product image and returns an error if it doesn't pass criteria
func (img *ProductImage) Validate(fh *multipart.FileHeader) *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if fh == nil {
		errs.Add(Invalid("image", l, msgValidateBrandLogo))
	}
	if fh != nil && fh.Size > FileUploadSizeLimit {
		errs.Add(Invalid("image", l, msgValidateBrandLogoSize))
	}

	if !errs.IsZero() {
		return NewValidationError("ProductImage", msgInvalidBrand, "", errs)
	}
	return nil
}

// Validate validates the product image patch and returns an error if it doesn't pass criteria
func (patch *ProductImagePatch) Validate(fh *multipart.FileHeader) *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if fh != nil && fh.Size > FileUploadSizeLimit {
		errs.Add(Invalid("image", l, msgValidateProductImageSize))
	}

	if !errs.IsZero() {
		return NewValidationError("ProductImagePatch", msgInvalidBrand, "", errs)
	}
	return nil
}

// Patch patches the product image
func (img *ProductImage) Patch(patch *ProductImagePatch) {
	if patch.URL != nil {
		img.URL = patch.URL
	}
}

// PreSave will fill timestamps
func (img *ProductImage) PreSave() {
	now := time.Now()
	img.CreatedAt = &now
	img.UpdatedAt = img.CreatedAt
}

// PreUpdate sets the update timestamp
func (img *ProductImage) PreUpdate() {
	now := time.Now()
	img.UpdatedAt = &now
}

// SetImageDetails sets the product image_url and public_id
func (img *ProductImage) SetImageDetails(details *gocloudinary.ResourceDetails) {
	img.URL = &details.SecureURL
	img.PublicID = &details.PublicID
}
