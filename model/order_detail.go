package model

// OrderDetail ties order with product items
type OrderDetail struct {
	OrderID      int64  `json:"order_id" db:"order_id"`
	ProductID    int64  `json:"product_id" db:"product_id"`
	Quantity     int    `json:"quantity" db:"quantity"`
	HistoryPrice int    `json:"history_price" db:"history_price"`
	HistorySKU   string `json:"history_sku" db:"history_sku"`
}

// OrderInfo returns the order details info with the product data
type OrderInfo struct {
	OrderDetail
	Product
}
