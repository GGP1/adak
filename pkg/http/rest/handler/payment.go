package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/GGP1/palo/pkg/shopping/payment"
	"github.com/GGP1/palo/pkg/storage/cache"
	"github.com/jinzhu/gorm"

	"github.com/stripe/stripe-go"
)

// CreatePayment creates a new payment intent.
func CreatePayment(db *gorm.DB, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var card model.Card
		var order ordering.Order

		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		r.Body.Close()

		order, ok := cache.Get("order")
		if !ok {
			response.Error(w, r, http.StatusNotFound, fmt.Errorf("order not found"))
			return
		}

		if order.Status == ordering.PaidState ||
			order.Status == ordering.ShippingState ||
			order.Status == ordering.ShippedState {
			response.Error(w, r, http.StatusBadRequest, fmt.Errorf("This order has already been paid"))
			return
		}

		intent, err := payment.CreateIntent(&order, card)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if intent.Status == stripe.PaymentIntentStatusCanceled {
			response.Error(w, r, http.StatusBadRequest, fmt.Errorf("Invalid payment intent status: %s", intent.Status))
			return
		}

		order.Status = ordering.PaidState

		err = db.Model(&order).Where("id=?", order.ID).UpdateColumn("status", order.Status).Error
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		respond := fmt.Sprintf("Thanks for your purchase! Your products will be delivered on %v.", order.DeliveryDate)
		fmt.Fprintln(w, respond)
	}
}
