package model

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"time"

	"github.com/dankobgd/ecommerce-shop/gocloudinary"
	"github.com/dankobgd/ecommerce-shop/utils/is"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/jmoiron/sqlx/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// error msgs
var (
	msgInvalidCategory            = &i18n.Message{ID: "model.category.validate.app_error", Other: "invalid category data"}
	msgValidateCategoryID         = &i18n.Message{ID: "model.category.validate.id.app_error", Other: "invalid category id"}
	msgValidateCategoryName       = &i18n.Message{ID: "model.category.validate.name.app_error", Other: "invalid category name"}
	msgValidateCategorySlug       = &i18n.Message{ID: "model.category.validate.created_at.app_error", Other: "invalid category slug"}
	msgValidateCategoryLogo       = &i18n.Message{ID: "model.category.validate.logo.app_error", Other: "invalid logo"}
	msgValidateCategoryLogoSize   = &i18n.Message{ID: "model.category.validate.logo.app_error", Other: "category file size exceeded, max 3MB"}
	msgValidateCategoryCrAt       = &i18n.Message{ID: "model.category.validate.created_at.app_error", Other: "invalid category created_at timestamp"}
	msgValidateCategoryUpAt       = &i18n.Message{ID: "model.category.validate.updated_at.app_error", Other: "invalid category updated_at timestamp"}
	msgValidateCategoryProperties = &i18n.Message{ID: "model.category.validate.properties.app_error", Other: "invalid json provided as properties"}
)

// Category is the category
type Category struct {
	TotalRecordsCount
	ID             int64           `json:"id" db:"id" schema:"-"`
	Name           string          `json:"name" db:"name" schema:"name"`
	Slug           string          `json:"slug" db:"slug" schema:"slug"`
	Description    string          `json:"description,omitempty" db:"description" schema:"description"`
	IsFeatured     bool            `json:"is_featured" db:"is_featured" schema:"is_featured"`
	Logo           string          `json:"logo" db:"logo" schema:"-"`
	LogoPublicID   string          `json:"logo_public_id" db:"logo_public_id" schema:"-"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at" schema:"-"`
	Properties     *types.JSONText `json:"properties" db:"properties" schema:"-"`
	PropertiesText *string         `json:"-" schema:"properties"`
}

// Validate validates the category and returns an error if it doesn't pass criteria
func (c *Category) Validate(fh *multipart.FileHeader) *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if c.ID != 0 {
		errs.Add(Invalid("id", l, msgValidateCategoryID))
	}
	if c.Name == "" {
		errs.Add(Invalid("name", l, msgValidateCategoryName))
	}
	if c.Slug == "" {
		errs.Add(Invalid("slug", l, msgValidateCategorySlug))
	}
	if c.CreatedAt.IsZero() {
		errs.Add(Invalid("created_at", l, msgValidateCategoryCrAt))
	}
	if c.UpdatedAt.IsZero() {
		errs.Add(Invalid("updated_at", l, msgValidateCategoryUpAt))
	}
	if fh == nil {
		errs.Add(Invalid("logo", l, msgValidateCategoryLogo))
	}
	if fh != nil && fh.Size > FileUploadSizeLimit {
		errs.Add(Invalid("logo", l, msgValidateCategoryLogoSize))
	}
	// ideally validate properties against json schema to check for the right keys, values and structure...
	if c.PropertiesText != nil && !is.ValidJSON(*c.PropertiesText) {
		errs.Add(Invalid("properties", l, msgValidateCategoryProperties))
	}

	if !errs.IsZero() {
		return NewValidationError("Category", msgInvalidCategory, "", errs)
	}
	return nil
}

// CategoryPatch is the category patch model
type CategoryPatch struct {
	Name           *string         `json:"name,omitempty" schema:"name"`
	Slug           *string         `json:"slug,omitempty" schema:"slug"`
	Description    *string         `json:"description,omitempty" schema:"description"`
	IsFeatured     *bool           `json:"is_featured,omitempty" schema:"is_featured"`
	Logo           *string         `json:"logo,omitempty" schema:"-"`
	Properties     *types.JSONText `json:"properties,omitempty" schema:"-"`
	PropertiesText *string         `json:"-" schema:"properties"`
}

// Patch patches the category fields that are provided
func (c *Category) Patch(patch *CategoryPatch) {
	if patch.Name != nil {
		c.Name = *patch.Name
	}
	if patch.Slug != nil {
		c.Slug = *patch.Slug
	}
	if patch.Description != nil {
		c.Description = *patch.Description
	}
	if patch.IsFeatured != nil {
		c.IsFeatured = *patch.IsFeatured
	}

	if patch.PropertiesText != nil {
		if *patch.PropertiesText == "" {
			c.Properties = nil
		} else {
			c.Properties = patch.Properties
		}
	}
}

// Validate validates the category patch and returns an error if it doesn't pass criteria
func (patch *CategoryPatch) Validate(fh *multipart.FileHeader) *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if patch.Name != nil && *patch.Name == "" {
		errs.Add(Invalid("name", l, msgValidateCategoryName))
	}
	if patch.Slug != nil && *patch.Slug == "" {
		errs.Add(Invalid("slug", l, msgValidateCategorySlug))
	}
	if fh != nil && fh.Size > FileUploadSizeLimit {
		errs.Add(Invalid("logo", l, msgValidateCategoryLogoSize))
	}
	// ideally validate properties against json schema to check for the right keys, values and structure...
	if patch.PropertiesText != nil && *patch.PropertiesText != "" && !is.ValidJSON(*patch.PropertiesText) {
		errs.Add(Invalid("properties", l, msgValidateCategoryProperties))
	}

	if !errs.IsZero() {
		return NewValidationError("Category", msgInvalidCategory, "", errs)
	}
	return nil
}

// CategoryPatchFromJSON decodes the input and returns the CategoryPatch
func CategoryPatchFromJSON(data io.Reader) (*CategoryPatch, error) {
	var patch *CategoryPatch
	err := json.NewDecoder(data).Decode(&patch)
	return patch, err
}

// CategoryFromJSON decodes the input and returns the Category
func CategoryFromJSON(data io.Reader) (*Category, error) {
	var c *Category
	err := json.NewDecoder(data).Decode(&c)
	return c, err
}

// ToJSON converts Category to json string
func (c *Category) ToJSON() string {
	b, _ := json.Marshal(c)
	return string(b)
}

// PreSave will set missing defaults and fill CreatedAt and UpdatedAt times
func (c *Category) PreSave() {
	c.CreatedAt = time.Now()
	c.UpdatedAt = c.CreatedAt
	c.SetProperties(c.PropertiesText)
}

// PreUpdate sets the update timestamp
func (c *Category) PreUpdate() {
	c.UpdatedAt = time.Now()
}

// SetProperties sets the category properties
func (c *Category) SetProperties(properties *string) {
	if properties != nil {
		props := types.JSONText(*properties)
		c.Properties = &props
	}
}

// SetProperties sets the category patch properties
func (patch *CategoryPatch) SetProperties(properties *string) {
	if properties != nil {
		if len(*properties) == 0 {
			patch.Properties = nil
		} else {
			props := types.JSONText(*properties)
			patch.Properties = &props
		}
	}
}

// SetLogoDetails sets the category logo and public_id
func (c *Category) SetLogoDetails(details *gocloudinary.ResourceDetails) {
	c.Logo = details.SecureURL
	c.LogoPublicID = details.PublicID
}
