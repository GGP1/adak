package stripe

import (
	"fmt"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/payout"
)

// CancelPayout cancels a payout. Funds will be refunded to the available balance.
func CancelPayout(payoutID string) (*stripe.Payout, error) {
	p, err := payout.Cancel(payoutID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe: Payout: %v", err)
	}

	return p, nil
}

// CreatePayout sends funds to the bank account.
func CreatePayout(amount int) (*stripe.Payout, error) {
	params := &stripe.PayoutParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}

	p, err := payout.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe: Payout: %v", err)
	}

	return p, nil
}

// GetPayout retrieves the details of an existing payout.
func GetPayout(payoutID string) (*stripe.Payout, error) {
	p, err := payout.Get(payoutID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe: Payout: %v", err)
	}

	return p, nil
}

// ListPayouts returns a list of existing payouts sent to third-party bank
// accounts or that Stripe has sent you.
func ListPayouts() []*stripe.Payout {
	var list []*stripe.Payout

	i := payout.List(nil)

	for i.Next() {
		p := i.Payout()
		list = append(list, p)
	}

	return list
}

// UpdatePayout updates the specified payout by setting the values of the
// parameters passed.
func UpdatePayout(payoutID string) (*stripe.Payout, error) {
	p, err := payout.Update(payoutID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe: Payout: %v", err)
	}

	return p, nil
}
