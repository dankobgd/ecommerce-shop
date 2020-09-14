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
	msgInvalidCategory             = &i18n.Message{ID: "model.category.validate.app_error", Other: "invalid category data"}
	msgValidateCategoryID          = &i18n.Message{ID: "model.category.validate.id.app_error", Other: "invalid category id"}
	msgValidateCategoryName        = &i18n.Message{ID: "model.category.validate.name.app_error", Other: "invalid category name"}
	msgValidateCategorySlug        = &i18n.Message{ID: "model.category.validate.created_at.app_error", Other: "invalid category slug"}
	msgValidateCategoryDescription = &i18n.Message{ID: "model.category.validate.updated_at.app_error", Other: "invalid category description"}
	msgValidateCategoryLogo        = &i18n.Message{ID: "model.category.validate.logo.app_error", Other: "invalid logo"}
	msgValidateCategoryCrAt        = &i18n.Message{ID: "model.category.validate.created_at.app_error", Other: "invalid category created_at timestamp"}
	msgValidateCategoryUpAt        = &i18n.Message{ID: "model.category.validate.updated_at.app_error", Other: "invalid category updated_at timestamp"}
)

// Category is the category
type Category struct {
	TotalRecordsCount
	ID          int64     `json:"id" db:"id" schema:"-"`
	Name        string    `json:"name" db:"name" schema:"name"`
	Slug        string    `json:"slug" db:"slug" schema:"slug"`
	Description string    `json:"description" db:"description" schema:"description"`
	IsFeatured  bool      `json:"is_featured" db:"is_featured" schema:"is_featured"`
	Logo        string    `json:"logo" db:"logo" schema:"-"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" schema:"-"`
}

// Validate validates the category and returns an error if it doesn't pass criteria
func (c *Category) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if c.ID != 0 {
		errs.Add(Invalid("category.id", l, msgValidateCategoryID))
	}
	if c.Name == "" {
		errs.Add(Invalid("category.name", l, msgValidateCategoryName))
	}
	if c.Slug == "" {
		errs.Add(Invalid("category.slug", l, msgValidateCategorySlug))
	}
	if c.Description == "" {
		errs.Add(Invalid("category.description", l, msgValidateCategoryDescription))
	}
	if c.CreatedAt.IsZero() {
		errs.Add(Invalid("category.created_at", l, msgValidateCategoryCrAt))
	}
	if c.UpdatedAt.IsZero() {
		errs.Add(Invalid("category.updated_at", l, msgValidateCategoryUpAt))
	}

	if !errs.IsZero() {
		return NewValidationError("Category", msgInvalidCategory, "", errs)
	}
	return nil
}

// CategoryPatch is the category patch model
type CategoryPatch struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	IsFeatured  *bool   `json:"is_featured"`
	Logo        *string `json:"logo"`
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
	if patch.Logo != nil {
		c.Logo = *patch.Logo
	}
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
}

// PreUpdate sets the update timestamp
func (c *Category) PreUpdate() {
	c.UpdatedAt = time.Now()
}

// SetLogoURL sets the logo url
func (c *Category) SetLogoURL(url string) {
	c.Logo = url
}
