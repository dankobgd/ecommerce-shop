package model

// OrderDetail ties order with product items
type OrderDetail struct {
	OrderID       int64  `json:"order_id" db:"order_id"`
	ProductID     int64  `json:"product_id" db:"product_id"`
	Quantity      int    `json:"quantity" db:"quantity"`
	OriginalPrice int    `json:"original_price" db:"original_price"`
	OriginalSKU   string `json:"original_sku" db:"original_sku"`
}
