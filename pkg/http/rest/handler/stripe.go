package handler

import (
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/shopping/payment/stripe"

	"github.com/go-chi/chi"
)

// StripeGetBalance responds with the account balance.
func StripeGetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		balance, err := stripe.GetBalance()
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		response.JSON(w, r, http.StatusOK, balance)
	}
}

// StripeGetEvent looks for the details of the events.
func StripeGetEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := chi.URLParam(r, "event")

		event, err := stripe.GetEvent(e)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		response.JSON(w, r, http.StatusOK, event)
	}
}

// StripeGetTxBalance responds with the transaction balance.
func StripeGetTxBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txID := chi.URLParam(r, "txID")

		tx, err := stripe.GetTxBalance(txID)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		response.JSON(w, r, http.StatusOK, tx)
	}
}

// StripeListEvents retrieves a list of all the stripe events within
// the last 30 days.
func StripeListEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventList := stripe.ListEvents()

		response.JSON(w, r, http.StatusOK, eventList)
	}
}

// StripeListTxs responds with a list of stripe transactions.
func StripeListTxs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txList := stripe.ListTxs()

		response.JSON(w, r, http.StatusOK, txList)
	}
}
