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
	"github.com/go-playground/validator"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// OrderParams holds the parameters for creating a order.
type OrderParams struct {
	Currency string      `json:"currency" validate:"required"`
	Address  string      `json:"address" validate:"required"`
	City     string      `json:"city" validate:"required"`
	Country  string      `json:"country" validate:"required"`
	State    string      `json:"state" validate:"required"`
	ZipCode  string      `json:"zip_code" validate:"required"`
	Date     date        `json:"date" validate:"required"`
	Card     stripe.Card `json:"card" validate:"required"`
}

type date struct {
	Year    int `json:"year" validate:"required,min=2020,max=2100"`
	Month   int `json:"month" validate:"required,min=1,max=12"`
	Day     int `json:"day" validate:"required,min=1,max=31"`
	Hour    int `json:"hour" validate:"required,min=0,max=24"`
	Minutes int `json:"minutes" validate:"required,min=0,max=60"`
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

		response.HTMLText(w, r, http.StatusOK, "The order has been deleted successfully.")
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

// GetByID retrieves all the orders from the user.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		order, err := GetByID(ctx, h.DB, id)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, order)
	}
}

// GetByUserID retrieves all the orders from the user.
func (h *Handler) GetByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uID, _ := r.Cookie("UID")

		ctx := r.Context()

		if err := token.CheckPermits(id, uID.Value); err != nil {
			response.Error(w, r, http.StatusUnauthorized, err)
			return
		}

		orders, err := GetByUserID(ctx, h.DB, id)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, orders)
	}
}

// New creates a new order and the payment intent.
func (h *Handler) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var oParams OrderParams
		cID, _ := r.Cookie("CID")
		uID, _ := r.Cookie("UID")
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&oParams); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
		}
		defer r.Body.Close()

		err := validator.New().StructCtx(ctx, oParams)
		if err != nil {
			http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
			return
		}

		// Parse jwt to take the user id
		userID, err := token.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Format date
		deliveryDate := time.Date(oParams.Date.Year, time.Month(oParams.Date.Month), oParams.Date.Day, oParams.Date.Hour, oParams.Date.Minutes, 0, 0, time.Local)

		if deliveryDate.Sub(time.Now()) < 0 {
			response.Error(w, r, http.StatusBadRequest, errors.New("past dates are not valid"))
			return
		}

		// Fetch the user cart
		cart, err := cart.Get(ctx, h.DB, cID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Create order passing userID, order params, delivery date and the user cart
		order, err := New(ctx, h.DB, userID, oParams, deliveryDate, cart)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Create payment intent and update the order status
		_, err = stripe.CreateIntent(order.ID, order.CartID, order.Currency, order.Cart.Total, oParams.Card)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := UpdateStatus(ctx, h.DB, order.ID, PaidState); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		respond := fmt.Sprintf("Thanks for your purchase! Your products will be delivered on %v.", order.DeliveryDate)
		response.HTMLText(w, r, http.StatusCreated, respond)
	}
}
