package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/ordering"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

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
		year := chi.URLParam(r, "year")
		month := chi.URLParam(r, "month")
		day := chi.URLParam(r, "day")
		hour := chi.URLParam(r, "hour")
		minutes := chi.URLParam(r, "minutes")

		y, _ := strconv.Atoi(year)
		m, _ := strconv.Atoi(month)
		d, _ := strconv.Atoi(day)
		h, _ := strconv.Atoi(hour)
		min, _ := strconv.Atoi(minutes)

		date := time.Date(y, time.Month(m), d, h, min, 0, 0, time.Local)

		if date.Sub(time.Now()) < 0 {
			response.Error(w, r, http.StatusBadRequest, errors.New("past dates are not valid"))
			return
		}

		cID, _ := r.Cookie("CID")
		uID, _ := r.Cookie("UID")

		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		cart, err := shopping.Get(db, cID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		order, err := ordering.New(db, userID.(string), cart, date)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		orderJSON, err := json.Marshal(order)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		_, err = http.Post("http://127.0.0.1:4000/payment", "application/json", bytes.NewBuffer(orderJSON))
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, order)
	}
}
