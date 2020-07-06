package model

// ProductImage is the product image
type ProductImage struct {
	ID        int64  `json:"id" db:"id"`
	ProductID int64  `json:"-" db:"product_id"`
	URL       string `json:"url" db:"url"`
}
