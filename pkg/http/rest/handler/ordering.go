package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/GGP1/palo/pkg/shopping/payment/stripe"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
)

// orderParams holds the parameters for creating a order
type orderParams struct {
	Currency string     `json:"currency"`
	Address  string     `json:"address"`
	City     string     `json:"city"`
	Country  string     `json:"country"`
	State    string     `json:"state"`
	ZipCode  string     `json:"zip_code"`
	Date     date       `json:"date"`
	Card     model.Card `json:"card"`
}

type date struct {
	Year    int `json:"year"`
	Month   int `json:"month"`
	Day     int `json:"day"`
	Hour    int `json:"hour"`
	Minutes int `json:"minutes"`
}

// DeleteOrder deletes an order.
func DeleteOrder(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := ordering.Delete(db, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "The order has been deleted.")
	}
}

// GetOrder finds all the stored orders.
func GetOrder(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var orders []ordering.Order

		if err := ordering.Get(db, &orders); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, orders)
	}
}

// NewOrder creates a new order.
func NewOrder(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cID, _ := r.Cookie("CID")
		uID, _ := r.Cookie("UID")

		var o orderParams

		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
		}
		defer r.Body.Close()

		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		deliveryDate := time.Date(o.Date.Year, time.Month(o.Date.Month), o.Date.Day, o.Date.Hour, o.Date.Minutes, 0, 0, time.Local)

		if deliveryDate.Sub(time.Now()) < 0 {
			response.Error(w, r, http.StatusBadRequest, fmt.Errorf("past dates are not valid"))
			return
		}

		cart, err := shopping.Get(db, cID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		order, err := ordering.NewOrder(db, userID, o.Currency, o.Address, o.City, o.Country, o.State, o.ZipCode, deliveryDate, cart)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		_, err = stripe.CreateIntent(order, o.Card)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		order.Status = ordering.PaidState

		_, err = db.Exec("UPDATE orders SET status=$2 WHERE id=$1", order.ID, order.Status)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		respond := fmt.Sprintf("Thanks for your purchase! Your products will be delivered on %v.", order.DeliveryDate)
		response.HTMLText(w, r, http.StatusOK, respond)
	}
}
