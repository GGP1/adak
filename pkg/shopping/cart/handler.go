package cart

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"

	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Handler manages cart endpoints.
type Handler struct {
	DB    *sqlx.DB
	Cache *lru.Cache
}

// Add appends a product to the cart.
func (h *Handler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Product
		q := chi.URLParam(r, "quantity")
		ctx := r.Context()

		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if quantity == 0 {
			response.Error(w, http.StatusBadRequest, errors.New("please insert a valid quantity"))
			return
		}

		if err = json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		cartID, _ := cookie.Get(r, "CID")

		cart, err := Add(ctx, h.DB, cartID.Value, &product, quantity)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, cart)
	}
}

// Checkout returns the final purchase.
func (h *Handler) Checkout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		checkout, err := Checkout(ctx, h.DB, cartID.Value)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, checkout)
	}
}

// FilterByBrand returns the products filtered by brand.
func (h *Handler) FilterByBrand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		brand := chi.URLParam(r, "brand")

		if err := sanitize.Normalize(&brand); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		products, err := FilterByBrand(ctx, h.DB, cartID.Value, brand)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// FilterByCategory returns the products filtered by category.
func (h *Handler) FilterByCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		category := chi.URLParam(r, "category")

		if err := sanitize.Normalize(&category); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		products, err := FilterByCategory(ctx, h.DB, cartID.Value, category)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// FilterByDiscount returns the products filtered by discount.
func (h *Handler) FilterByDiscount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")

		minDiscount, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		maxDiscount, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		products, err := FilterByDiscount(ctx, h.DB, cartID.Value, minDiscount, maxDiscount)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// FilterBySubtotal returns the products filtered by subtotal.
func (h *Handler) FilterBySubtotal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")

		minSubtotal, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		maxSubtotal, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		products, err := FilterBySubtotal(ctx, h.DB, cartID.Value, minSubtotal, maxSubtotal)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// FilterByTaxes returns the products filtered by taxes.
func (h *Handler) FilterByTaxes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")

		minTaxes, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		maxTaxes, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		products, err := FilterByTaxes(ctx, h.DB, cartID.Value, minTaxes, maxTaxes)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// FilterByTotal returns the products filtered by total.
func (h *Handler) FilterByTotal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")

		minTotal, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		maxTotal, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		products, err := FilterByTotal(ctx, h.DB, cartID.Value, minTotal, maxTotal)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// FilterByType returns the products filtered by type.
func (h *Handler) FilterByType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		pType := chi.URLParam(r, "type")

		if err := sanitize.Normalize(&pType); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		products, err := FilterByType(ctx, h.DB, cartID.Value, pType)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// FilterByWeight returns the products filtered by weight.
func (h *Handler) FilterByWeight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		ctx := r.Context()

		minWeight, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		maxWeight, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		products, err := FilterByWeight(ctx, h.DB, cartID.Value, minWeight, maxWeight)
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
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if cCart, ok := h.Cache.Get(cartID); ok {
			response.JSON(w, http.StatusOK, cCart)
			return
		}

		cart, err := Get(ctx, h.DB, cartID.Value)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		h.Cache.Add(cartID, cart)
		response.JSON(w, http.StatusOK, cart)
	}
}

// Products retrieves cart products.
func (h *Handler) Products() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		items, err := Products(ctx, h.DB, cartID.Value)
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
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		id := chi.URLParam(r, "id")
		q := chi.URLParam(r, "quantity")
		ctx := r.Context()

		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if err := Remove(ctx, h.DB, cartID.Value, id, quantity); err != nil {
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
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if err := Reset(ctx, h.DB, cartID.Value); err != nil {
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
		cartID, err := cookie.Get(r, "CID")
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		size, err := Size(ctx, h.DB, cartID.Value)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, size)
	}
}
