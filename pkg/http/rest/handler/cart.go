package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
)

// CartAdd appends a product to the cart
func CartAdd(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := chi.URLParam(r, "quantity")

		var (
			product shopping.CartProduct
			ctx     = r.Context()
		)

		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		if quantity == 0 {
			response.Error(w, r, http.StatusBadRequest, fmt.Errorf("error: please insert a valid quantity"))
			return
		}

		if err = json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		cartID, _ := r.Cookie("CID")

		cart, err := shopping.Add(ctx, db, cartID.Value, &product, quantity)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusCreated, cart)
	}
}

// CartCheckout returns the final purchase
func CartCheckout(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		checkout, err := shopping.Checkout(ctx, db, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, checkout)
	}
}

// CartFilterByBrand returns the products filtered by brand
func CartFilterByBrand(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		brand := chi.URLParam(r, "brand")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		products, err := shopping.FilterByBrand(ctx, db, c.Value, brand)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByCategory returns the products filtered by category
func CartFilterByCategory(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := chi.URLParam(r, "category")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		products, err := shopping.FilterByCategory(ctx, db, c.Value, category)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByDiscount returns the products filtered by discount
func CartFilterByDiscount(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minDiscount, _ := strconv.ParseFloat(min, 64)
		maxDiscount, _ := strconv.ParseFloat(max, 64)

		products, err := shopping.FilterByDiscount(ctx, db, c.Value, minDiscount, maxDiscount)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterBySubtotal returns the products filtered by subtotal
func CartFilterBySubtotal(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minSubtotal, _ := strconv.ParseFloat(min, 64)
		maxSubtotal, _ := strconv.ParseFloat(max, 64)

		products, err := shopping.FilterBySubtotal(ctx, db, c.Value, minSubtotal, maxSubtotal)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByTaxes returns the products filtered by taxes
func CartFilterByTaxes(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minTaxes, _ := strconv.ParseFloat(min, 64)
		maxTaxes, _ := strconv.ParseFloat(max, 64)

		products, err := shopping.FilterByTaxes(ctx, db, c.Value, minTaxes, maxTaxes)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByTotal returns the products filtered by total
func CartFilterByTotal(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minTotal, _ := strconv.ParseFloat(min, 64)
		maxTotal, _ := strconv.ParseFloat(max, 64)

		products, err := shopping.FilterByTotal(ctx, db, c.Value, minTotal, maxTotal)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByType returns the products filtered by type
func CartFilterByType(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productType := chi.URLParam(r, "type")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		products, err := shopping.FilterByType(ctx, db, c.Value, productType)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByWeight returns the products filtered by weight
func CartFilterByWeight(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		minWeight, _ := strconv.ParseFloat(min, 64)
		maxWeight, _ := strconv.ParseFloat(max, 64)

		products, err := shopping.FilterByWeight(ctx, db, c.Value, minWeight, maxWeight)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartGet returns the cart in a JSON format
func CartGet(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		cart, err := shopping.Get(ctx, db, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, cart)
	}
}

// CartItems prints cart items
func CartItems(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		items, err := shopping.Items(ctx, db, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, items)
	}
}

// CartRemove takes out a product from the shopping cart
func CartRemove(db *sqlx.DB) http.HandlerFunc {
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

		if err = shopping.Remove(ctx, db, c.Value, id, quantity); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Successfully removed the product from the cart")
	}
}

// CartReset resets the cart to its default state
func CartReset(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		if err := shopping.Reset(ctx, db, c.Value); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Cart reseted")
	}
}

// CartSize returns the size of the shopping cart
func CartSize(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		ctx := r.Context()

		size, err := shopping.Size(ctx, db, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, size)
	}
}
