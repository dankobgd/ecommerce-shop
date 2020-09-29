package payment

import "github.com/dankobgd/ecommerce-shop/model"

// Provider is the payment processor service
type Provider interface {
	Name() string
	Charge(paymentID string, order *model.Order, user *model.User, amount uint64, currency string) (string, error)
	Refund(paymentID string, amount uint64, currency string) (string, error)
	Confirm(paymentID string) error
}
