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
	msgInvalidTag = &i18n.Message{ID: "model.tag.validate.app_error", Other: "invalid tag data"}

	msgValidateTagID          = &i18n.Message{ID: "model.tag.validate.id.app_error", Other: "invalid  tag id"}
	msgValidateTagProductID   = &i18n.Message{ID: "model.tag.validate.product_id.app_error", Other: "invalid tag product id"}
	msgValidateTagName        = &i18n.Message{ID: "model.tag.validate.name.app_error", Other: "invalid tag name"}
	msgValidateTagSlug        = &i18n.Message{ID: "model.tag.validate.slug.app_error", Other: "invalid tag slug"}
	msgValidateTagDescription = &i18n.Message{ID: "model.tag.validate.description.app_error", Other: "invalid tag description"}
	msgValidateTagCrAt        = &i18n.Message{ID: "model.tag.validate.created_at.app_error", Other: "invalid tag created_at timestamp"}
	msgValidateTagUpAt        = &i18n.Message{ID: "model.tag.validate.updated_at.app_error", Other: "invalid tag updated_at timestamp"}
)

// Tag is the tag model
type Tag struct {
	TotalRecordsCount
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

// TagPatch is the patch for tag
type TagPatch struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
}

// Patch patches the product tag
func (t *Tag) Patch(patch *TagPatch) {
	if patch.Name != nil {
		t.Name = *patch.Name
	}
	if patch.Slug != nil {
		t.Slug = *patch.Slug
	}
	if patch.Description != nil {
		t.Description = *patch.Description
	}
}

// TagFromJSON decodes the input and returns the Tag
func TagFromJSON(data io.Reader) (*Tag, error) {
	var t *Tag
	err := json.NewDecoder(data).Decode(&t)
	return t, err
}

// TagPatchFromJSON decodes the input and returns the TagPatch
func TagPatchFromJSON(data io.Reader) (*TagPatch, error) {
	var p *TagPatch
	err := json.NewDecoder(data).Decode(&p)
	return p, err
}

// PreSave will fill timestamps
func (t *Tag) PreSave() {
	t.CreatedAt = time.Now()
	t.UpdatedAt = t.CreatedAt
}

// PreUpdate sets the update timestamp
func (t *Tag) PreUpdate() {
	t.UpdatedAt = time.Now()
}

// Validate validates the tag and returns an error if it doesn't pass criteria
func (t *Tag) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if t.ID != 0 {
		errs.Add(Invalid("id", l, msgValidateTagID))
	}
	if t.Name == "" {
		errs.Add(Invalid("name", l, msgValidateTagName))
	}
	if t.Slug == "" {
		errs.Add(Invalid("slug", l, msgValidateTagSlug))
	}
	if t.Description == "" {
		errs.Add(Invalid("description", l, msgValidateTagDescription))
	}
	if t.CreatedAt.IsZero() {
		errs.Add(Invalid("created_at", l, msgValidateTagCrAt))
	}
	if t.UpdatedAt.IsZero() {
		errs.Add(Invalid("updated_at", l, msgValidateTagUpAt))
	}

	if !errs.IsZero() {
		return NewValidationError("Tag", msgInvalidTag, "", errs)
	}
	return nil
}
