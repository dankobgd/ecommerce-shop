package model

// ProductImage is the product image
type ProductImage struct {
	ID        int64  `json:"id" db:"img_id"`
	ProductID int64  `json:"-" db:"img_product_id"`
	URL       string `json:"url" db:"img_url"`
}
