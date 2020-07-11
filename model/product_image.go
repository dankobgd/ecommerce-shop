package model

import "time"

// ProductImage is the product image
type ProductImage struct {
	ID        *int64     `json:"id" db:"img_id"`
	ProductID *int64     `json:"-" db:"img_product_id"`
	URL       *string    `json:"url" db:"img_url"`
	CreatedAt *time.Time `json:"created_at" db:"img_created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"img_updated_at"`
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
