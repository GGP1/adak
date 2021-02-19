package stripe

import (
	"github.com/pkg/errors"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

// CancelIntent cancels the purchase.
func CancelIntent(intentID string) error {
	_, err := paymentintent.Cancel(intentID, nil)
	if err != nil {
		return errors.Wrap(err, "stripe: PaymentIntent")
	}

	return nil
}

// CaptureIntent captures the funds of an existing uncaptured PaymentIntent
//  when its status is requires_capture.
func CaptureIntent(intentID string) error {
	_, err := paymentintent.Capture(intentID, nil)
	if err != nil {
		return errors.Wrap(err, "stripe: PaymentIntent")
	}

	return nil
}

// ConfirmIntent confirms that your customer intends to pay with current
// or provided payment method.
func ConfirmIntent(intentID string, source *stripe.Source) error {
	pi, err := paymentintent.Get(intentID, nil)
	if err != nil {
		return errors.Wrap(err, "stripe: PaymentIntent")
	}

	if pi.Status != "requires_payment_method" {
		return errors.Errorf("stripe: paymentIntent already has a status of %s", pi.Status)
	}

	params := &stripe.PaymentIntentConfirmParams{
		Source: stripe.String(source.ID),
	}

	_, err = paymentintent.Confirm(pi.ID, params)
	if err != nil {
		return errors.Wrap(err, "stripe: PaymentIntent")
	}

	return nil
}

// CreateIntent creates a payment intent object.
func CreateIntent(id, cartID, currency string, total int64, card Card) (*stripe.PaymentIntent, error) {
	pMethodID, err := CreateMethod(card)
	if err != nil {
		return nil, err
	}

	if total < 50 {
		return nil, errors.New("stripe: the order total should be higher than $0.50")
	}

	// Amounts to be provided in a currency’s smallest unit
	// 100 = 1 USD
	// minimum: $0.50 / maximum: $999,999.99
	params := &stripe.PaymentIntentParams{
		PaymentMethod: stripe.String(pMethodID),
		Amount:        stripe.Int64(total),
		Currency:      stripe.String(currency),
		ConfirmationMethod: stripe.String(string(
			stripe.PaymentIntentConfirmationMethodManual,
		)),
		Confirm: stripe.Bool(true),
		Params: stripe.Params{
			Metadata: map[string]string{
				"order_id": id,
				"cart_id":  cartID,
			},
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, errors.Wrap(err, "stripe: PaymentIntent")
	}

	if pi.Status == stripe.PaymentIntentStatusCanceled {
		return nil, errors.Wrapf(err, "stripe: invalid PaymentIntent status: %s", pi.Status)
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
		return nil, errors.Wrap(err, "stripe: PaymentIntent")
	}

	return pi, nil
}

// UpdateIntent sets new properties on a PaymentIntent object without confirming.
func UpdateIntent(intentID string, total int64) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(total),
	}

	pi, err := paymentintent.Update(intentID, params)
	if err != nil {
		return nil, errors.Wrap(err, "stripe: PaymentIntent")
	}

	return pi, nil
}
