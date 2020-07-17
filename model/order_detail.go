package model

import (
	"encoding/json"
	"io"
)

// OrderDetail ties order with product items
type OrderDetail struct {
	OrderID       int64  `json:"order_id" db:"order_id"`
	ProductID     int64  `json:"product_id" db:"product_id"`
	Quantity      int    `json:"quantity" db:"quantity"`
	OriginalPrice int    `json:"original_price" db:"original_price"`
	OriginalSKU   string `json:"original_sku" db:"original_sku"`
}

// OrderItemData is the request data
type OrderItemData struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// OrderItemsDataFromJSON decodes the input and returns the order item data list
func OrderItemsDataFromJSON(data io.Reader) ([]*OrderItemData, error) {
	var od []*OrderItemData
	err := json.NewDecoder(data).Decode(&od)
	return od, err
}
