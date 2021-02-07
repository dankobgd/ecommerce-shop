package stripe

import (
	"errors"
	"fmt"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/payment"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type stripePaymentProvider struct {
	client *client.API
}

// NewPaymentProvider returns the stripe payment provider
func NewPaymentProvider(secretKey string) (payment.Provider, *model.AppErr) {
	s := &stripePaymentProvider{
		client: &client.API{},
	}

	s.client.Init(secretKey, nil)

	return s, nil
}

func (sp *stripePaymentProvider) Name() string {
	return "Stripe"
}

func (sp *stripePaymentProvider) Charge(paymentID string, order *model.Order, user *model.User, amount uint64, currency string) (*stripe.PaymentIntent, error) {
	return sp.chargePaymentIntent(paymentID, amount, currency, order, user)
}

func (sp *stripePaymentProvider) Refund(transactionID string, amount uint64, currency string) (string, error) {
	stripeAmount := int64(amount)
	ref, err := sp.client.Refunds.New(&stripe.RefundParams{
		Charge: &transactionID,
		Amount: &stripeAmount,
	})
	if err != nil {
		return "", err
	}

	return ref.ID, err
}

func (sp *stripePaymentProvider) Confirm(paymentID string) error {
	_, err := sp.client.PaymentIntents.Confirm(paymentID, nil)

	if stripeErr, ok := err.(*stripe.Error); ok {
		return errors.New(stripeErr.Msg)
	}

	return err
}

func (sp *stripePaymentProvider) chargePaymentIntent(paymentMethodID string, amount uint64, currency string, order *model.Order, user *model.User) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		PaymentMethod: stripe.String(paymentMethodID),
		Amount:        stripe.Int64(int64(amount)),
		Currency:      stripe.String(currency),
		Shipping:      prepareShippingAddress(order, user),
		Confirm:       stripe.Bool(true),
	}
	intent, err := sp.client.PaymentIntents.New(params)
	if err != nil {
		return nil, err
	}

	if intent.Status == stripe.PaymentIntentStatusRequiresAction {
		return intent, fmt.Errorf("PaymentIntent status %s", intent.ClientSecret)
	}

	if intent.Status == stripe.PaymentIntentStatusSucceeded {
		return intent, nil
	}

	return nil, fmt.Errorf("Invalid PaymentIntent status: %s", intent.Status)
}

func prepareShippingAddress(order *model.Order, user *model.User) *stripe.ShippingDetailsParams {
	return &stripe.ShippingDetailsParams{
		Address: &stripe.AddressParams{
			Line1:      &order.ShippingAddressLine1,
			Line2:      order.ShippingAddressLine2,
			City:       &order.ShippingAddressCity,
			State:      order.ShippingAddressState,
			Country:    &order.ShippingAddressCountry,
			PostalCode: order.ShippingAddressZIP,
		},
		Name: stripe.String(fmt.Sprintf("%s %s", user.FirstName, user.LastName)),
	}
}
