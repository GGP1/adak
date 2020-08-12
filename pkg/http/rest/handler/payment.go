package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/GGP1/palo/pkg/shopping/payment"

	"github.com/stripe/stripe-go"
)

// CreatePayment creates a new payment intent.
func CreatePayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stripe.Key = cfg.StripeKey
		var order ordering.Order

		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		r.Body.Close()

		if order.Status == ordering.PaidState ||
			order.Status == ordering.ShippingState ||
			order.Status == ordering.ShippedState {
			response.Error(w, r, http.StatusBadRequest, fmt.Errorf("This order has already been paid"))
			return
		}

		clientSecret, err := payment.CreateIntent(&order)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, struct {
			ClientSecret string `json:"clientSecret"`
		}{
			ClientSecret: clientSecret,
		})
	}
}
