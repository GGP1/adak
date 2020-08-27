package ordering

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/internal/token"
	"github.com/GGP1/palo/pkg/shopping/cart"
	"github.com/GGP1/palo/pkg/shopping/payment/stripe"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// OrderParams holds the parameters for creating a order.
type OrderParams struct {
	Currency string      `json:"currency"`
	Address  string      `json:"address"`
	City     string      `json:"city"`
	Country  string      `json:"country"`
	State    string      `json:"state"`
	ZipCode  string      `json:"zip_code"`
	Date     date        `json:"date"`
	Card     stripe.Card `json:"card"`
}

type date struct {
	Year    int `json:"year"`
	Month   int `json:"month"`
	Day     int `json:"day"`
	Hour    int `json:"hour"`
	Minutes int `json:"minutes"`
}

// Handler handles ordering endpoints.
type Handler struct {
	DB *sqlx.DB
}

// Delete deletes an order.
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		if err := Delete(ctx, h.DB, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "The order has been deleted succesfully.")
	}
}

// Get finds all the stored orders.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orders, err := Get(ctx, h.DB)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, orders)
	}
}

// New creates a new order.
func (h *Handler) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cID, _ := r.Cookie("CID")
		uID, _ := r.Cookie("UID")

		var (
			o   OrderParams
			ctx = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
		}
		defer r.Body.Close()

		err := o.validate()
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		userID, err := token.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		deliveryDate := time.Date(o.Date.Year, time.Month(o.Date.Month), o.Date.Day, o.Date.Hour, o.Date.Minutes, 0, 0, time.Local)

		if deliveryDate.Sub(time.Now()) < 0 {
			response.Error(w, r, http.StatusBadRequest, errors.New("past dates are not valid"))
			return
		}

		cart, err := cart.Get(ctx, h.DB, cID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		order, err := New(ctx, h.DB, userID, o.Currency, o.Address, o.City, o.Country, o.State, o.ZipCode, deliveryDate, cart)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		_, err = stripe.CreateIntent(order.ID, order.CartID, order.Currency, order.Cart.Total, o.Card)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		order.Status = PaidState

		_, err = h.DB.ExecContext(ctx, "UPDATE orders SET status=$2 WHERE id=$1", order.ID, order.Status)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		respond := fmt.Sprintf("Thanks for your purchase! Your products will be delivered on %v.", order.DeliveryDate)
		response.HTMLText(w, r, http.StatusCreated, respond)
	}
}

// validate order input.
func (o *OrderParams) validate() error {
	if o.Address == "" {
		return errors.New("address is required")
	}

	if o.Currency == "" {
		return errors.New("currency is required")
	}

	if o.City == "" {
		return errors.New("city is required")
	}

	if o.Country == "" {
		return errors.New("country is required")
	}

	if o.State == "" {
		return errors.New("state is required")
	}

	if o.ZipCode == "" {
		return errors.New("zipcode is required")
	}

	if o.Card.Number == "" {
		return errors.New("card number is required")
	}

	if o.Card.ExpMonth == "" || len(o.Card.ExpMonth) > 2 {
		return errors.New("card expiration month is required")
	}

	if o.Card.ExpYear == "" || len(o.Card.ExpYear) > 5 {
		return errors.New("invalid card expiration year")
	}

	if o.Card.CVC == "" || len(o.Card.CVC) > 3 {
		return errors.New("invalid card cvc")
	}

	if o.Date.Year < 2020 || o.Date.Year > 2100 {
		return errors.New("invalid year")
	}

	if o.Date.Month < 1 || o.Date.Month > 12 {
		return errors.New("invalid month")
	}

	if o.Date.Day < 1 || o.Date.Day > 31 {
		return errors.New("invalid day")
	}

	if o.Date.Hour < 1 || o.Date.Hour > 24 {
		return errors.New("invalid hour")
	}

	if o.Date.Minutes < 1 || o.Date.Minutes > 60 {
		return errors.New("invalid minutes")
	}

	return nil
}
