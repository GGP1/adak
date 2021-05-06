package stripe

import (
	"github.com/pkg/errors"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/refund"
)

// CreateRefund will refund a charge that has previously been created but not yet
// refunded.
// Funds will be refunded to the credit or debit card that was originally charged.
func CreateRefund(intentID string) (*stripe.Refund, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(intentID),
	}

	r, err := refund.New(params)
	if err != nil {
		return nil, errors.Wrap(err, "stripe: Refund")
	}

	return r, nil
}

// GetRefund retrieves the details of an existing refund.
func GetRefund(refundID string) (*stripe.Refund, error) {
	r, err := refund.Get(refundID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "stripe: Refund")
	}

	return r, nil
}

// ListRefunds returns a list of all refunds you’ve previously created.
func ListRefunds() []*stripe.Refund {
	var list []*stripe.Refund

	i := refund.List(nil)

	for i.Next() {
		bt := i.Refund()
		list = append(list, bt)
	}

	return list
}

// UpdateRefund updates the specified refund by setting the values of the parameters
// passed. Any parameters not provided will be left unchanged.
func UpdateRefund(refundID string) (*stripe.Refund, error) {
	r, err := refund.Update(refundID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "stripe: Refund")
	}

	return r, nil
}
