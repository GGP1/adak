package payment

import (
	"fmt"

	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping/ordering"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/paymentmethod"
)

// CreatePaymentMethod creates a new payment method.
func CreatePaymentMethod(card model.Card) (string, error) {
	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(card.Number),
			ExpMonth: stripe.String(card.ExpMonth),
			ExpYear:  stripe.String(card.ExpYear),
			CVC:      stripe.String(card.CVC),
		},
		Type: stripe.String("card"),
	}

	pm, err := paymentmethod.New(params)
	if err != nil {
		return "", fmt.Errorf("Invalid card parameters")
	}

	return pm.ID, nil
}

// CreateIntent charges the purchase.
func CreateIntent(order *ordering.Order, card model.Card) (*stripe.PaymentIntent, error) {
	pMethodID, err := CreatePaymentMethod(card)
	if err != nil {
		return nil, err
	}

	// Amounts to be provided in a currency’s smallest unit
	// 100 = 1 USD
	// minimum: $0.50 / maximum: $999,999.99
	amount := order.Cart.Total * 100

	if amount < 100 {
		return nil, fmt.Errorf("the order total should be higher than $1")
	}

	params := &stripe.PaymentIntentParams{
		PaymentMethod: stripe.String(pMethodID),
		Amount:        stripe.Int64(int64(amount)),
		Currency:      stripe.String(order.Currency),
		ConfirmationMethod: stripe.String(string(
			stripe.PaymentIntentConfirmationMethodManual,
		)),
		Confirm: stripe.Bool(true),
		Params: stripe.Params{
			Metadata: map[string]string{
				"order_id": order.ID,
				"cart_id":  order.CartID,
			},
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("payments: error creating payment intent: %v", err)
	}

	if pi.Status == stripe.PaymentIntentStatusCanceled {
		return nil, fmt.Errorf("Invalid PaymentIntent status: %s", pi.Status)
	}

	return pi, nil
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
