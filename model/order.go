package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var msgInvalidOrderData = &i18n.Message{ID: "model.order.validate.app_error", Other: "Invalid order data"}
var msgValidatePaymentMethodID = &i18n.Message{ID: "model.order.validate.payment_method_id.app_error", Other: "Payment method id is required"}
var msgValidateNoItems = &i18n.Message{ID: "model.order.validate.no_items.app_error", Other: "No order items provided"}
var msgValidateBillingAddress = &i18n.Message{ID: "model.order.validate.billing_address.app_error", Other: "Invalid billing address"}
var msgValidateBillingAddressID = &i18n.Message{ID: "model.order.validate.billing_address_id.app_error", Other: "Invalid billing address id"}
var msgValidateShippingAddress = &i18n.Message{ID: "model.order.validate.shipping_address.app_error", Other: "Invalid shipping address"}
var msgValidateShippingAddressNeedsBilling = &i18n.Message{ID: "model.order.validate.shipping_address.app_error", Other: "No billing address provided but same_shipping_as_billing is true"}

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
	PromoCode                *string    `json:"promo_code" db:"promo_code"`
	Status                   string     `json:"status" db:"status"`
	Subtotal                 int        `json:"subtotal" db:"subtotal"`
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
	PaymentMethodID           string      `json:"payment_method_id"`
	Items                     []*CartItem `json:"items"`
	BillingAddress            *Address    `json:"billing_address"`
	ShippingAddress           *Address    `json:"shipping_address"`
	SaveAddress               *bool       `json:"save_address"`
	UseExistingBillingAddress *bool       `json:"use_existing_billing_address"`
	BillingAddressID          *int64      `json:"billing_address_id"`
	SameShippingAsBilling     *bool       `json:"same_shipping_as_billing"`
	PromoCode                 *string     `json:"promo_code"`
}

// OrderRequestDataFromJSON decodes the input and returns the order item data list
func OrderRequestDataFromJSON(data io.Reader) (*OrderRequestData, error) {
	var ord *OrderRequestData
	err := json.NewDecoder(data).Decode(&ord)
	return ord, err
}

// Validate validates the tag and returns an error if it doesn't pass criteria
func (data *OrderRequestData) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if data.PaymentMethodID == "" {
		errs.Add(Invalid("payment_method_id", l, msgValidatePaymentMethodID))
	}
	if len(data.Items) == 0 {
		errs.Add(Invalid("items", l, msgValidateNoItems))
	}

	if data.BillingAddress == nil && (data.UseExistingBillingAddress == nil || (data.UseExistingBillingAddress != nil && *data.UseExistingBillingAddress == false)) {
		errs.Add(Invalid("billing_address", l, msgValidateBillingAddress))
	}
	if data.BillingAddressID == nil && data.UseExistingBillingAddress != nil && *data.UseExistingBillingAddress == true {
		errs.Add(Invalid("billing_address_id", l, msgValidateBillingAddressID))
	}
	if (data.ShippingAddress == nil && (data.UseExistingBillingAddress == nil || (data.UseExistingBillingAddress != nil && *data.UseExistingBillingAddress == false))) && (data.SameShippingAsBilling == nil || (data.SameShippingAsBilling != nil && *data.SameShippingAsBilling == false)) {
		errs.Add(Invalid("shipping_address", l, msgValidateShippingAddress))
	}
	if (data.BillingAddress != nil && data.BillingAddressID != nil) && data.SameShippingAsBilling != nil && *data.SameShippingAsBilling == true {
		errs.Add(Invalid("shipping_address", l, msgValidateShippingAddressNeedsBilling))
	}

	if !errs.IsZero() {
		return NewValidationError("OrderRequestData", msgInvalidOrderData, "", errs)
	}
	return nil
}
