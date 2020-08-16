package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/GGP1/palo/pkg/shopping/payment/stripe"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
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
func DeleteOrder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		orderID, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = ordering.Delete(db, orderID)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "The order has been deleted.")
	}
}

// GetOrder finds all the stored orders.
func GetOrder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var orders []ordering.Order

		err := ordering.Get(db, &orders)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, orders)
	}
}

// NewOrder creates a new order.
func NewOrder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cID, _ := r.Cookie("CID")
		uID, _ := r.Cookie("UID")

		var o orderParams

		err := json.NewDecoder(r.Body).Decode(&o)
		if err != nil {
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

		err = db.Model(&order).Where("id=?", order.ID).UpdateColumn("status", order.Status).Error
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		respond := fmt.Sprintf("Thanks for your purchase! Your products will be delivered on %v.", order.DeliveryDate)
		response.HTMLText(w, r, http.StatusOK, respond)
	}
}
