package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

// orderParams holds the parameters for creating a order
type orderParams struct {
	Currency string `json:"currency"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Country  string `json:"country"`
	State    string `json:"state"`
	ZipCode  string `json:"zip_code"`
	Date     date   `json:"date"`
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
			response.Error(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		deliveryDate := time.Date(o.Date.Year, time.Month(o.Date.Month), o.Date.Day, o.Date.Hour, o.Date.Minutes, 0, 0, time.Local)

		if deliveryDate.Sub(time.Now()) < 0 {
			response.Error(w, r, http.StatusBadRequest, errors.New("past dates are not valid"))
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

		orderJSON, err := json.Marshal(order)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		go func() {
			_, err = http.Post("http://127.0.0.1:4000/payment", "application/json", bytes.NewBuffer(orderJSON))
			if err != nil {
				response.Error(w, r, http.StatusInternalServerError, err)
				return
			}
		}()

		respond := fmt.Sprintf("Thanks for your purchase! Your products will be delivered on %v.", order.DeliveryDate)

		response.JSON(w, r, http.StatusOK, order)
		fmt.Fprintln(w, respond)
	}
}
