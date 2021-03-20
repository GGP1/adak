package ordering

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/shopping/payment/stripe"

	"github.com/go-chi/chi"
	validator "github.com/go-playground/validator/v10"
	lru "github.com/hashicorp/golang-lru"
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
	Year    int `json:"year" validate:"required,min=2021,max=2150"`
	Month   int `json:"month" validate:"required,min=1,max=12"`
	Day     int `json:"day" validate:"required,min=1,max=31"`
	Hour    int `json:"hour" validate:"required,min=0,max=24"`
	Minutes int `json:"minutes" validate:"required,min=0,max=60"`
}

// Handler handles ordering endpoints.
type Handler struct {
	DB    *sqlx.DB
	Cache *lru.Cache
}

// Delete deletes an order.
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := Delete(ctx, h.DB, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, fmt.Sprintf("order %q deleted", id))
	}
}

// Get finds all the stored orders.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orders, err := Get(ctx, h.DB)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, orders)
	}
}

// GetByID retrieves all the orders from the user.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		order, err := GetByID(ctx, h.DB, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, order)
	}
}

// GetByUserID retrieves all the orders from the user.
func (h *Handler) GetByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()
		userID, err := cookie.Get(r, "UID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		if err := token.CheckPermits(id, userID.Value); err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		// Distinguish from the other ids from the same user
		cacheKey := fmt.Sprintf("%s orders", id)
		if cOrders, ok := h.Cache.Get(cacheKey); ok {
			response.JSON(w, http.StatusOK, cOrders)
			return
		}

		orders, err := GetByUserID(ctx, h.DB, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		h.Cache.Add(cacheKey, orders)
		response.JSON(w, http.StatusOK, orders)
	}
}

// New creates a new order and the payment intent.
func (h *Handler) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var oParams OrderParams
		ctx := r.Context()
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}
		uID, err := cookie.Get(r, "UID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&oParams); err != nil {
			response.Error(w, http.StatusBadRequest, err)
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, oParams); err != nil {
			response.Error(w, http.StatusBadRequest, err.(validator.ValidationErrors))
			return
		}

		if err := sanitize.Normalize(&oParams.Address, &oParams.City, &oParams.Country, &oParams.Currency, &oParams.State, &oParams.ZipCode); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		userID, err := token.GetUserID(uID.Value)
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		// Format date
		deliveryDate := time.Date(oParams.Date.Year, time.Month(oParams.Date.Month), oParams.Date.Day, oParams.Date.Hour, oParams.Date.Minutes, 0, 0, time.Local)
		if deliveryDate.Sub(time.Now()) < 0 {
			response.Error(w, http.StatusBadRequest, errors.New("past dates are not valid"))
			return
		}

		// Fetch the user cart
		cart, err := cart.Get(ctx, h.DB, cartID.Value)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		// Create order passing userID, order params, delivery date and the user cart
		order, err := New(ctx, h.DB, userID, oParams, deliveryDate, cart)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		// Create payment intent and update the order status
		_, err = stripe.CreateIntent(order.ID, order.CartID, order.Currency, order.Cart.Total, oParams.Card)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if err := UpdateStatus(ctx, h.DB, order.ID, PaidState); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusCreated, fmt.Sprintf("order %q created", order.ID))
	}
}
