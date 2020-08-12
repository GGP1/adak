package payment

import (
	"fmt"

	"github.com/GGP1/palo/pkg/shopping/ordering"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

// CreateIntent charges the purchase.
func CreateIntent(order *ordering.Order) (string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(order.Cart.Total) * 10),
		Currency: stripe.String(order.Currency),
		Params: stripe.Params{
			Metadata: map[string]string{
				"order_id": order.ID,
				"cart_id":  order.CartID,
			},
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", fmt.Errorf("payments: error creating payment intent: %v", err)
	}

	if pi.Status == stripe.PaymentIntentStatusCanceled {
		return "", fmt.Errorf("Invalid PaymentIntent status: %s", pi.Status)
	}

	return pi.ClientSecret, nil
}

// RetrieveIntent retires the purchase.
func RetrieveIntent(paymentIntent string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Get(paymentIntent, nil)
	if err != nil {
		return nil, fmt.Errorf("payments: error fetching payment intent: %v", err)
	}

	return pi, nil
}

// ConfirmIntent confirms the purchase.
func ConfirmIntent(paymentIntent string, source *stripe.Source) error {
	pi, err := paymentintent.Get(paymentIntent, nil)
	if err != nil {
		return fmt.Errorf("error fetching payment intent for confirmation: %v", err)
	}

	if pi.Status != "requires_payment_method" {
		return fmt.Errorf("paymentIntent already has a status of %s", pi.Status)
	}

	params := &stripe.PaymentIntentConfirmParams{
		Source: stripe.String(source.ID),
	}

	_, err = paymentintent.Confirm(pi.ID, params)
	if err != nil {
		return fmt.Errorf("error confirming PaymentIntent: %v", err)
	}

	return nil
}

// CancelIntent cancels the purchase.
func CancelIntent(paymentIntent string) error {
	_, err := paymentintent.Cancel(paymentIntent, nil)
	if err != nil {
		return fmt.Errorf("payments: error canceling PaymentIntent: %v", err)
	}

	return nil
}

// UpdateIntent sets intent new values.
func UpdateIntent(paymentIntent string, order ordering.Order) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(int64(order.Cart.Total)),
	}

	pi, err := paymentintent.Update(paymentIntent, params)
	if err != nil {
		return nil, fmt.Errorf("payments: error updating payment intent: %v", err)
	}

	return pi, nil
}
