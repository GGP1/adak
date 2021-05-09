package ordering

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/internal/validate"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/shopping/payment/stripe"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4/zero"
)

// OrderParams holds the parameters for creating a order.
type OrderParams struct {
	Currency string      `json:"currency" validate:"required"`
	Address  string      `json:"address" validate:"required"`
	City     string      `json:"city" validate:"required"`
	Country  string      `json:"country" validate:"required"`
	State    string      `json:"state" validate:"required"`
	ZipCode  string      `json:"zip_code" validate:"required"`
	Date     Date        `json:"date" validate:"required"`
	Card     stripe.Card `json:"card" validate:"required"`
}

// Date of the order.
type Date struct {
	Year    int `json:"year" validate:"required,min=2021,max=2150"`
	Month   int `json:"month" validate:"required,min=1,max=12"`
	Day     int `json:"day" validate:"required,min=1,max=31"`
	Hour    int `json:"hour" validate:"required,min=0,max=24"`
	Minutes int `json:"minutes" validate:"required,min=0,max=60"`
}

// Handler handles ordering endpoints.
type Handler struct {
	Development     bool
	OrderingService Service
	CartService     cart.Service
	DB              *sqlx.DB
	Cache           *memcache.Client
}

// Delete deletes an order.
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := h.OrderingService.Delete(ctx, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, id)
	}
}

// Get finds all the stored orders.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orders, err := h.OrderingService.Get(ctx)
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

		order, err := h.OrderingService.GetByID(ctx, id)
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

		if err := token.CheckPermits(r, id); err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		// Every service has ids of different length, they will never match
		item, err := h.Cache.Get(id)
		if err == nil {
			response.EncodedJSON(w, item.Value)
			return
		}

		orders, err := h.OrderingService.GetByUserID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSONAndCache(h.Cache, w, id, orders)
	}
}

// New creates a new order and the payment intent.
func (h *Handler) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}
		userID, err := cookie.GetValue(r, "UID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		var orderParams OrderParams
		if err := json.NewDecoder(r.Body).Decode(&orderParams); err != nil {
			response.Error(w, http.StatusBadRequest, err)
		}
		defer r.Body.Close()

		if err := validateOrderParams(ctx, &orderParams); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		id := token.RandString(30)
		order, err := h.OrderingService.New(ctx, id, userID, cartID, orderParams, h.CartService)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if !h.Development {
			// Create payment intent and update the order status
			_, err = stripe.CreateIntent(order.ID.String, order.CartID.String,
				order.Currency.String, order.Cart.Total.Int64, orderParams.Card)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, err)
				return
			}
		}

		if err := h.OrderingService.UpdateStatus(ctx, order.ID.String, zero.IntFrom(int64(Paid))); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if err := h.CartService.Reset(ctx, cartID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, order)
	}
}

func validateOrderParams(ctx context.Context, oParams *OrderParams) error {
	if err := validate.Struct(ctx, oParams); err != nil {
		return err
	}
	oParams.Address = sanitize.Normalize(oParams.Address)
	oParams.City = sanitize.Normalize(oParams.City)
	oParams.Country = sanitize.Normalize(oParams.Country)
	oParams.Currency = sanitize.Normalize(oParams.Currency)
	oParams.State = sanitize.Normalize(oParams.State)
	oParams.ZipCode = sanitize.Normalize(oParams.ZipCode)

	return nil
}
