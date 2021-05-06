package stripe

import (
	"github.com/pkg/errors"

	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/balance"
	"github.com/stripe/stripe-go/v72/balancetransaction"
)

// GetBalance retrieves the current account balance,
// based on the authentication that was used to make the request.
func GetBalance() (*stripe.Balance, error) {
	b, err := balance.Get(nil)
	if err != nil {
		return nil, errors.Wrap(err, "stripe: account balance")
	}

	return b, nil
}

// GetTxBalance retrieves the balance transaction with the given ID.
func GetTxBalance(txID string) (*stripe.BalanceTransaction, error) {
	txBalance, err := balancetransaction.Get(txID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "stripe: transaction balance")
	}

	return txBalance, nil
}

// ListTxs returns a list of transactions that have contributed to the
// Stripe account balance.
func ListTxs() []*stripe.BalanceTransaction {
	var list []*stripe.BalanceTransaction

	i := balancetransaction.List(nil)

	for i.Next() {
		bt := i.BalanceTransaction()
		list = append(list, bt)
	}

	return list
}
