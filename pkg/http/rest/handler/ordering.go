package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/GGP1/palo/pkg/shopping/ordering"

	"github.com/GGP1/palo/pkg/shopping/payment/stripe"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
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

		response.HTMLText(w, r, http.StatusOK, "The order has been deleted succesfully.")
	}
}

// GetOrder finds all the stored orders.
func GetOrder(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		orders, err := ordering.Get(db)
		if err != nil {
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

		err := o.validate()
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
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
		response.HTMLText(w, r, http.StatusCreated, respond)
	}
}

// validate order input.
func (o *orderParams) validate() error {
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
