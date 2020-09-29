package model

import (
	"encoding/json"
	"io"
	"time"
)

type orderStatus int

// order statuses
const (
	OrderStatusPending orderStatus = iota
	OrderStatusSuccess
	OrderStatusFailed
)

func (s orderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "pending"
	case OrderStatusSuccess:
		return "success"
	case OrderStatusFailed:
		return "fail"
	default:
		return "unknown"
	}
}

// Order represents the transaction
type Order struct {
	TotalRecordsCount
	ID                       int64      `json:"id" db:"id"`
	UserID                   int64      `json:"user_id" db:"user_id"`
	Status                   string     `json:"status" db:"status"`
	Total                    int        `json:"total" db:"total"`
	ShippedAt                *time.Time `json:"shipped_at" db:"shipped_at"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	BillingAddressLine1      string     `json:"billing_address_line_1,omitempty" db:"billing_address_line_1"`
	BillingAddressLine2      *string    `json:"billing_address_line_2,omitempty" db:"billing_address_line_2"`
	BillingAddressCity       string     `json:"billing_address_city,omitempty" db:"billing_address_city"`
	BillingAddressCountry    string     `json:"billing_address_country,omitempty" db:"billing_address_country"`
	BillingAddressState      *string    `json:"billing_address_state,omitempty" db:"billing_address_state"`
	BillingAddressZIP        *string    `json:"billing_address_zip,omitempty" db:"billing_address_zip"`
	BillingAddressLatitude   *float64   `json:"billing_address_latitude,omitempty" db:"billing_address_latitude"`
	BillingAddressLongitude  *float64   `json:"billing_address_longitude,omitempty" db:"billing_address_longitude"`
	ShippingAddressLine1     string     `json:"shipping_address_line_1,omitempty" db:"shipping_address_line_1"`
	ShippingAddressLine2     *string    `json:"shipping_address_line_2,omitempty" db:"shipping_address_line_2"`
	ShippingAddressCity      string     `json:"shipping_address_city,omitempty" db:"shipping_address_city"`
	ShippingAddressCountry   string     `json:"shipping_address_country,omitempty" db:"shipping_address_country"`
	ShippingAddressState     *string    `json:"shipping_address_state,omitempty" db:"shipping_address_state"`
	ShippingAddressZIP       *string    `json:"shipping_address_zip,omitempty" db:"shipping_address_zip"`
	ShippingAddressLatitude  *float64   `json:"shipping_address_latitude,omitempty" db:"shipping_address_latitude"`
	ShippingAddressLongitude *float64   `json:"shipping_address_longitude,omitempty" db:"shipping_address_longitude"`
}

// PreSave fills the defaults
func (o *Order) PreSave() {
	o.CreatedAt = time.Now()
	if o.Status == "" {
		o.Status = OrderStatusPending.String()
	}
}

// CartItem is the cart item info
type CartItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// OrderRequestData is used to create new order
type OrderRequestData struct {
	PaymentMethodID       string      `json:"payment_method_id"`
	SameShippingAsBilling bool        `json:"same_shipping_as_billing"`
	Items                 []*CartItem `json:"items"`
	BillingAddress        *Address    `json:"billing_address"`
	ShippingAddress       *Address    `json:"shipping_address"`
}

// OrderRequestDataFromJSON decodes the input and returns the order item data list
func OrderRequestDataFromJSON(data io.Reader) (*OrderRequestData, error) {
	var ord *OrderRequestData
	err := json.NewDecoder(data).Decode(&ord)
	return ord, err
}
