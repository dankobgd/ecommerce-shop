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
	msgInvalidBrand             = &i18n.Message{ID: "model.brand.validate.app_error", Other: "invalid brand data"}
	msgValidateBrandID          = &i18n.Message{ID: "model.brand.validate.id.app_error", Other: "invalid brand id"}
	msgValidateBrandName        = &i18n.Message{ID: "model.brand.validate.name.app_error", Other: "invalid brand name"}
	msgValidateBrandSlug        = &i18n.Message{ID: "model.brand.validate.slug.app_error", Other: "invalid brand slug"}
	msgValidateBrandDescription = &i18n.Message{ID: "model.brand.validate.description.app_error", Other: "invalid brand description"}
	msgValidateBrandType        = &i18n.Message{ID: "model.brand.validate.type.app_error", Other: "invalid brand type"}
	msgValidateBrandEmail       = &i18n.Message{ID: "model.brand.validate.email.app_error", Other: "invalid brand email"}
	msgValidateBrandLogo        = &i18n.Message{ID: "model.brand.validate.logo.app_error", Other: "invalid brand logo"}
	msgValidateBrandWebsiteURL  = &i18n.Message{ID: "model.brand.validate.website_url.app_error", Other: "invalid brand website URL"}
	msgValidateBrandCrAt        = &i18n.Message{ID: "model.brand.validate.created_at.app_error", Other: "invalid brand created_at timestamp"}
	msgValidateBrandUpAt        = &i18n.Message{ID: "model.brand.validate.updated_at.app_error", Other: "invalid brand updated_at timestamp"}
)

// Brand is the brand of the product
type Brand struct {
	ID          int64     `json:"id" db:"id" schema:"-"`
	Name        string    `json:"name" db:"name" schema:"name"`
	Slug        string    `json:"slug" db:"slug" schema:"slug"`
	Type        string    `json:"type" db:"type" schema:"type"`
	Description string    `json:"description" db:"description" schema:"description"`
	Email       string    `json:"email" db:"email" schema:"email"`
	WebsiteURL  string    `json:"website_url" db:"website_url" schema:"website_url"`
	Logo        string    `json:"logo" db:"logo" schema:"-"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" schema:"-"`
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
func (b *Brand) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if b.ID != 0 {
		errs.Add(Invalid("brand.id", l, msgValidateBrandID))
	}
	if b.Name == "" {
		errs.Add(Invalid("brand.name", l, msgValidateBrandName))
	}
	if b.Slug == "" {
		errs.Add(Invalid("brand.slug", l, msgValidateBrandSlug))
	}
	if b.Type == "" {
		errs.Add(Invalid("brand.type", l, msgValidateBrandType))
	}
	if b.Description == "" {
		errs.Add(Invalid("brand.description", l, msgValidateBrandDescription))
	}
	if b.Email == "" {
		errs.Add(Invalid("brand.email", l, msgValidateBrandEmail))
	}
	if b.WebsiteURL == "" {
		errs.Add(Invalid("brand.website_url", l, msgValidateBrandWebsiteURL))
	}
	if b.CreatedAt.IsZero() {
		errs.Add(Invalid("brand.created_at", l, msgValidateBrandCrAt))
	}
	if b.UpdatedAt.IsZero() {
		errs.Add(Invalid("brand.updated_at", l, msgValidateBrandUpAt))
	}

	if !errs.IsZero() {
		return NewValidationError("Brand", msgInvalidBrand, "", errs)
	}
	return nil
}

// BrandPatch is the brand patch model
type BrandPatch struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
	Email       *string `json:"email"`
	Logo        *string `json:"logo"`
	WebsiteURL  *string `json:"website_url"`
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
	if patch.Logo != nil {
		b.Logo = *patch.Logo
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

// SetLogoURL sets the logo url
func (b *Brand) SetLogoURL(url string) {
	b.Logo = url
}
