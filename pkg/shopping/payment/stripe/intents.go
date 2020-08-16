package stripe

import (
	"fmt"

	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping/ordering"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

// CancelIntent cancels the purchase.
func CancelIntent(intentID string) error {
	_, err := paymentintent.Cancel(intentID, nil)
	if err != nil {
		return fmt.Errorf("payments: error canceling PaymentIntent: %v", err)
	}

	return nil
}

// CaptureIntent captures the funds of an existing uncaptured PaymentIntent
//  when its status is requires_capture.
func CaptureIntent(intentID string) error {
	_, err := paymentintent.Capture(intentID, nil)
	if err != nil {
		return fmt.Errorf("payments: error capturing PaymentIntent: %v", err)
	}

	return nil
}

// ConfirmIntent confirms that your customer intends to pay with current
// or provided payment method.
func ConfirmIntent(intentID string, source *stripe.Source) error {
	pi, err := paymentintent.Get(intentID, nil)
	if err != nil {
		return fmt.Errorf("payments: error fetching PaymentIntent for confirmation: %v", err)
	}

	if pi.Status != "requires_payment_method" {
		return fmt.Errorf("payments: paymentIntent already has a status of %s", pi.Status)
	}

	params := &stripe.PaymentIntentConfirmParams{
		Source: stripe.String(source.ID),
	}

	_, err = paymentintent.Confirm(pi.ID, params)
	if err != nil {
		return fmt.Errorf("payments: error confirming PaymentIntent: %v", err)
	}

	return nil
}

// CreateIntent creates a payment intent object.
func CreateIntent(order *ordering.Order, card model.Card) (*stripe.PaymentIntent, error) {
	pMethodID, err := CreateMethod(card)
	if err != nil {
		return nil, err
	}

	// Amounts to be provided in a currency’s smallest unit
	// 100 = 1 USD
	// minimum: $0.50 / maximum: $999,999.99
	amount := order.Cart.Total * 100

	if amount < 100 {
		return nil, fmt.Errorf("payments: the order total should be higher than $1")
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
		return nil, fmt.Errorf("payments: error creating PaymentIntent: %v", err)
	}

	if pi.Status == stripe.PaymentIntentStatusCanceled {
		return nil, fmt.Errorf("payments: invalid PaymentIntent status: %s", pi.Status)
	}

	return pi, nil
}

// ListIntents returns a list of PaymentIntents.
func ListIntents() []*stripe.PaymentIntent {
	var list []*stripe.PaymentIntent

	params := &stripe.PaymentIntentListParams{}
	params.Filters.AddFilter("limit", "", "3")

	i := paymentintent.List(params)

	for i.Next() {
		list = append(list, i.PaymentIntent())
	}

	return list
}

// RetrieveIntent lists the details of a PaymentIntent that has previously been created.
func RetrieveIntent(intentID string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Get(intentID, nil)
	if err != nil {
		return nil, fmt.Errorf("payments: error fetching PaymentIntent: %v", err)
	}

	return pi, nil
}

// UpdateIntent sets new properties on a PaymentIntent object without confirming.
func UpdateIntent(intentID string, order ordering.Order) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(int64(order.Cart.Total)),
	}

	pi, err := paymentintent.Update(intentID, params)
	if err != nil {
		return nil, fmt.Errorf("payments: error updating PaymentIntent: %v", err)
	}

	return pi, nil
}
