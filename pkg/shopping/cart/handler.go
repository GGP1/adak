package cart

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GGP1/palo/internal/response"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Handler manages cart endpoints.
type Handler struct {
	DB *sqlx.DB
}

// Add appends a product to the cart.
func (h *Handler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := chi.URLParam(r, "quantity")

		var (
			product Product
			ctx     = r.Context()
		)

		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		if quantity == 0 {
			response.Error(w, r, http.StatusBadRequest, errors.New("please insert a valid quantity"))
			return
		}

		if err = json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		cartID, _ := r.Cookie("CID")

		cart, err := Add(ctx, h.DB, cartID.Value, &product, quantity)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusCreated, cart)
	}
}

// Checkout returns the final purchase.
func (h *Handler) Checkout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		checkout, err := Checkout(ctx, h.DB, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, checkout)
	}
}

// FilterByBrand returns the products filtered by brand.
func (h *Handler) FilterByBrand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		brand := chi.URLParam(r, "brand")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		products, err := FilterByBrand(ctx, h.DB, c.Value, brand)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// FilterByCategory returns the products filtered by category.
func (h *Handler) FilterByCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := chi.URLParam(r, "category")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		products, err := FilterByCategory(ctx, h.DB, c.Value, category)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// FilterByDiscount returns the products filtered by discount.
func (h *Handler) FilterByDiscount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minDiscount, _ := strconv.ParseFloat(min, 64)
		maxDiscount, _ := strconv.ParseFloat(max, 64)

		products, err := FilterByDiscount(ctx, h.DB, c.Value, minDiscount, maxDiscount)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// FilterBySubtotal returns the products filtered by subtotal.
func (h *Handler) FilterBySubtotal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minSubtotal, _ := strconv.ParseFloat(min, 64)
		maxSubtotal, _ := strconv.ParseFloat(max, 64)

		products, err := FilterBySubtotal(ctx, h.DB, c.Value, minSubtotal, maxSubtotal)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// FilterByTaxes returns the products filtered by taxes.
func (h *Handler) FilterByTaxes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minTaxes, _ := strconv.ParseFloat(min, 64)
		maxTaxes, _ := strconv.ParseFloat(max, 64)

		products, err := FilterByTaxes(ctx, h.DB, c.Value, minTaxes, maxTaxes)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// FilterByTotal returns the products filtered by total.
func (h *Handler) FilterByTotal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minTotal, _ := strconv.ParseFloat(min, 64)
		maxTotal, _ := strconv.ParseFloat(max, 64)

		products, err := FilterByTotal(ctx, h.DB, c.Value, minTotal, maxTotal)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// FilterByType returns the products filtered by type.
func (h *Handler) FilterByType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productType := chi.URLParam(r, "type")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		products, err := FilterByType(ctx, h.DB, c.Value, productType)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// FilterByWeight returns the products filtered by weight.
func (h *Handler) FilterByWeight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minWeight, _ := strconv.ParseFloat(min, 64)
		maxWeight, _ := strconv.ParseFloat(max, 64)

		products, err := FilterByWeight(ctx, h.DB, c.Value, minWeight, maxWeight)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// Get returns the cart in a JSON format.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		cart, err := Get(ctx, h.DB, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, cart)
	}
}

// Products retrieves cart products.
func (h *Handler) Products() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		items, err := Products(ctx, h.DB, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, items)
	}
}

// Remove takes out a product from the shopping cart.
func (h *Handler) Remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		q := chi.URLParam(r, "quantity")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = Remove(ctx, h.DB, c.Value, id, quantity); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Successfully removed the product from the cart")
	}
}

// Reset resets the cart to its default state.
func (h *Handler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		if err := Reset(ctx, h.DB, c.Value); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Cart reseted")
	}
}

// Size returns the size of the shopping cart.
func (h *Handler) Size() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		size, err := Size(ctx, h.DB, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, size)
	}
}
