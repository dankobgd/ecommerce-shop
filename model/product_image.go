package model

import (
	"encoding/json"
	"io"
	"time"
)

// ProductImage is the product image
type ProductImage struct {
	ID        *int64     `json:"id" db:"img_id"`
	ProductID *int64     `json:"-" db:"img_product_id"`
	URL       *string    `json:"url" db:"img_url"`
	CreatedAt *time.Time `json:"created_at" db:"img_created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"img_updated_at"`
}

// ProductImagePatch is the patch for tag
type ProductImagePatch struct {
	URL *string `json:"url"`
}

// ProductImagePatchFromJSON decodes the input and returns the ProductImagePatch
func ProductImagePatchFromJSON(data io.Reader) (*ProductImagePatch, error) {
	var p *ProductImagePatch
	err := json.NewDecoder(data).Decode(&p)
	return p, err
}

// Patch patches the product image
func (img *ProductImage) Patch(patch *ProductImagePatch) {
	if img.URL != nil {
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
