package cart

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/params"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/validate"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4/zero"
)

// Handler manages cart endpoints.
type Handler struct {
	service Service
	db      *sqlx.DB
	cache   *memcache.Client
}

// NewHandler returns a new cart handler.
func NewHandler(service Service, db *sqlx.DB, cache *memcache.Client) Handler {
	return Handler{
		service: service,
		db:      db,
		cache:   cache,
	}
}

// Add appends a product to the cart.
func (h *Handler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		var product Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validate.Struct(ctx, product); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		product.CartID = zero.StringFrom(cartID)
		if err := h.service.Add(ctx, product); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, product)
	}
}

// Checkout returns the final purchase.
func (h *Handler) Checkout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		checkout, err := h.service.Checkout(ctx, cartID)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, checkout)
	}
}

// FilterBy returns the products filtered by the field provided.
func (h *Handler) FilterBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		ctx := r.Context()
		field := sanitize.Normalize(chi.URLParam(r, "field"))
		args := sanitize.Normalize(chi.URLParam(r, "args"))

		products, err := h.service.FilterBy(ctx, cartID, field, args)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// Get returns the cart in a JSON format.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		item, err := h.cache.Get(cartID)
		if err == nil {
			response.EncodedJSON(w, item.Value)
			return
		}

		cart, err := h.service.Get(ctx, cartID)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSONAndCache(h.cache, w, cartID, cart)
	}
}

// Products retrieves cart products.
func (h *Handler) Products() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		items, err := h.service.CartProducts(ctx, cartID)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, items)
	}
}

// Remove takes out a product from the shopping cart.
func (h *Handler) Remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		ctx := r.Context()
		id, err := params.URLID(ctx)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		q := chi.URLParam(r, "quantity")
		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if err := h.service.Remove(ctx, cartID, id, int64(quantity)); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, fmt.Sprintf("product %q deleted from cart %q", id, cartID))
	}
}

// Reset resets the cart to its default state.
func (h *Handler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		if err := h.service.Reset(ctx, cartID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, fmt.Sprintf("cart %q reseted", cartID))
	}
}

// Size returns the size of the shopping cart.
func (h *Handler) Size() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		size, err := h.service.Size(ctx, cartID)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, size)
	}
}
