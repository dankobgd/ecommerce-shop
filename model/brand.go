package model

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"time"

	"github.com/dankobgd/ecommerce-shop/gocloudinary"
	"github.com/dankobgd/ecommerce-shop/utils/is"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// error msgs
var (
	msgInvalidBrand            = &i18n.Message{ID: "model.brand.validate.app_error", Other: "invalid brand data"}
	msgValidateBrandID         = &i18n.Message{ID: "model.brand.validate.id.app_error", Other: "invalid brand id"}
	msgValidateBrandName       = &i18n.Message{ID: "model.brand.validate.name.app_error", Other: "invalid brand name"}
	msgValidateBrandSlug       = &i18n.Message{ID: "model.brand.validate.slug.app_error", Other: "invalid brand slug"}
	msgValidateBrandType       = &i18n.Message{ID: "model.brand.validate.type.app_error", Other: "invalid brand type"}
	msgValidateBrandEmail      = &i18n.Message{ID: "model.brand.validate.email.app_error", Other: "invalid brand email"}
	msgValidateBrandWebsiteURL = &i18n.Message{ID: "model.brand.validate.website_url.app_error", Other: "invalid brand website URL"}
	msgValidateBrandCrAt       = &i18n.Message{ID: "model.brand.validate.created_at.app_error", Other: "invalid brand created_at timestamp"}
	msgValidateBrandUpAt       = &i18n.Message{ID: "model.brand.validate.updated_at.app_error", Other: "invalid brand updated_at timestamp"}
	msgValidateBrandLogo       = &i18n.Message{ID: "model.brand.validate.logo.app_error", Other: "invalid brand logo"}
	msgValidateBrandLogoSize   = &i18n.Message{ID: "model.brand.validate.logo.app_error", Other: "File size exceeded, max 3MB allowed"}
)

// Brand is the brand of the product
type Brand struct {
	TotalRecordsCount
	ID           int64     `json:"id" db:"id" schema:"-"`
	Name         string    `json:"name" db:"name" schema:"name"`
	Slug         string    `json:"slug" db:"slug" schema:"slug"`
	Type         string    `json:"type" db:"type" schema:"type"`
	Description  string    `json:"description,omitempty" db:"description" schema:"description"`
	Email        string    `json:"email" db:"email" schema:"email"`
	WebsiteURL   string    `json:"website_url" db:"website_url" schema:"website_url"`
	Logo         string    `json:"logo" db:"logo" schema:"-"`
	LogoPublicID string    `json:"logo_public_id" db:"logo_public_id" schema:"-"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" schema:"-"`
}

// PreSave will fill timestamps and other defaults
func (b *Brand) PreSave() {
	b.CreatedAt = time.Now()
	b.UpdatedAt = b.CreatedAt
	b.Email = NormalizeEmail(b.Email)
}

// PreUpdate sets the update timestamp
func (b *Brand) PreUpdate() {
	b.UpdatedAt = time.Now()
	b.Email = NormalizeEmail(b.Email)
}

// Validate validates the brand and returns an error if it doesn't pass criteria
func (b *Brand) Validate(fh *multipart.FileHeader) *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if b.ID != 0 {
		errs.Add(Invalid("id", l, msgValidateBrandID))
	}
	if b.Name == "" {
		errs.Add(Invalid("name", l, msgValidateBrandName))
	}
	if b.Slug == "" {
		errs.Add(Invalid("slug", l, msgValidateBrandSlug))
	}
	if b.Type == "" {
		errs.Add(Invalid("type", l, msgValidateBrandType))
	}
	if len(b.Email) == 0 || len(b.Email) > userEmailMaxLength || !is.ValidEmail(b.Email) {
		errs.Add(Invalid("email", l, msgValidateBrandEmail))
	}
	if b.WebsiteURL == "" {
		errs.Add(Invalid("website_url", l, msgValidateBrandWebsiteURL))
	}
	if b.CreatedAt.IsZero() {
		errs.Add(Invalid("created_at", l, msgValidateBrandCrAt))
	}
	if b.UpdatedAt.IsZero() {
		errs.Add(Invalid("updated_at", l, msgValidateBrandUpAt))
	}
	if fh == nil {
		errs.Add(Invalid("logo", l, msgValidateBrandLogo))
	}
	if fh != nil && fh.Size > FileUploadSizeLimit {
		errs.Add(Invalid("logo", l, msgValidateBrandLogoSize))
	}

	if !errs.IsZero() {
		return NewValidationError("Brand", msgInvalidBrand, "", errs)
	}
	return nil
}

// BrandPatch is the brand patch model
type BrandPatch struct {
	Name        *string `json:"name,omitempty" schema:"name"`
	Slug        *string `json:"slug,omitempty" schema:"slug"`
	Type        *string `json:"type,omitempty" schema:"type"`
	Description *string `json:"description,omitempty" schema:"description"`
	Email       *string `json:"email,omitempty" schema:"email"`
	Logo        *string `json:"logo,omitempty" schema:"-"`
	WebsiteURL  *string `json:"website_url,omitempty" schema:"website_url"`
}

// Patch patches the brand fields that are provided
func (b *Brand) Patch(patch *BrandPatch) {
	if patch.Name != nil {
		b.Name = *patch.Name
	}
	if patch.Slug != nil {
		b.Slug = *patch.Slug
	}
	if patch.Type != nil {
		b.Type = *patch.Type
	}
	if patch.Description != nil {
		b.Description = *patch.Description
	}
	if patch.Email != nil {
		b.Email = *patch.Email
	}
	if patch.WebsiteURL != nil {
		b.WebsiteURL = *patch.WebsiteURL
	}
}

// Validate validates the brand patch and returns an error if it doesn't pass criteria
func (patch *BrandPatch) Validate(fh *multipart.FileHeader) *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if patch.Name != nil && *patch.Name == "" {
		errs.Add(Invalid("name", l, msgValidateBrandName))
	}
	if patch.Slug != nil && *patch.Slug == "" {
		errs.Add(Invalid("slug", l, msgValidateBrandSlug))
	}
	if patch.Type != nil && *patch.Type == "" {
		errs.Add(Invalid("type", l, msgValidateBrandType))
	}
	if patch.Email != nil && (len(*patch.Email) == 0 || len(*patch.Email) > userEmailMaxLength || !is.ValidEmail(*patch.Email)) {
		errs.Add(Invalid("email", l, msgValidateBrandEmail))
	}
	if patch.WebsiteURL != nil && *patch.WebsiteURL == "" {
		errs.Add(Invalid("website_url", l, msgValidateBrandWebsiteURL))
	}
	if fh != nil && fh.Size > FileUploadSizeLimit {
		errs.Add(Invalid("logo", l, msgValidateBrandLogoSize))
	}

	if !errs.IsZero() {
		return NewValidationError("Brand", msgInvalidUser, "", errs)
	}
	return nil
}

// BrandPatchFromJSON decodes the input and returns the BrandPatch
func BrandPatchFromJSON(data io.Reader) (*BrandPatch, error) {
	var patch *BrandPatch
	err := json.NewDecoder(data).Decode(&patch)
	return patch, err
}

// BrandFromJSON decodes the input and returns the Brand
func BrandFromJSON(data io.Reader) (*Brand, error) {
	var b *Brand
	err := json.NewDecoder(data).Decode(&b)
	return b, err
}

// ToJSON converts Brand to json string
func (b *Brand) ToJSON() string {
	bytes, _ := json.Marshal(b)
	return string(bytes)
}

// SetLogoDetails sets the brand logo and public_id
func (b *Brand) SetLogoDetails(details *gocloudinary.ResourceDetails) {
	b.Logo = details.SecureURL
	b.LogoPublicID = details.PublicID
}
