package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/ordering"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

// CreatePayment creates a new payment intent
func CreatePayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stripe.Key = cfg.StripeKey
		var order ordering.Order

		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		r.Body.Close()

		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(10000),
			Currency: stripe.String(string(stripe.CurrencyUSD)),
			PaymentMethodTypes: stripe.StringSlice([]string{
				"card",
			}),
		}

		params.AddMetadata("ordered_at", order.OrderedAt.Format("15:04:05 02/01/2006"))
		params.AddMetadata("order_id", string(order.ID))

		pi, err := paymentintent.New(params)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, struct {
			ClientSecret string `json:"clientSecret"`
		}{
			ClientSecret: pi.ClientSecret,
		})
	}
}
