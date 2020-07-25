package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// CartAdd appends a product to the cart
func CartAdd(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product shopping.CartProduct

		q := mux.Vars(r)["quantity"]
		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		if quantity == 0 {
			response.Error(w, r, http.StatusBadRequest, fmt.Errorf("error: please insert a valid quantity"))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		c, _ := r.Cookie("CID")

		cart, err := shopping.Add(db, c.Value, &product, quantity)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, cart)
	}
}

// CartCheckout returns the final purchase
func CartCheckout(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")
		checkout, err := shopping.Checkout(db, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, checkout)
	}
}

// CartFilterByBrand returns the products filtered by brand
func CartFilterByBrand(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		brand := mux.Vars(r)["brand"]
		c, _ := r.Cookie("CID")

		products, err := shopping.FilterByBrand(db, c.Value, brand)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByCategory returns the products filtered by category
func CartFilterByCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := mux.Vars(r)["category"]
		c, _ := r.Cookie("CID")

		products, err := shopping.FilterByCategory(db, c.Value, category)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByDiscount returns the products filtered by discount
func CartFilterByDiscount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]
		c, _ := r.Cookie("CID")

		minDiscount, _ := strconv.ParseFloat(min, 32)
		maxDiscount, _ := strconv.ParseFloat(max, 32)

		products, err := shopping.FilterByDiscount(db, c.Value, float32(minDiscount), float32(maxDiscount))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterBySubtotal returns the products filtered by subtotal
func CartFilterBySubtotal(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]
		c, _ := r.Cookie("CID")

		minSubtotal, _ := strconv.ParseFloat(min, 32)
		maxSubtotal, _ := strconv.ParseFloat(max, 32)

		products, err := shopping.FilterBySubtotal(db, c.Value, float32(minSubtotal), float32(maxSubtotal))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByTaxes returns the products filtered by taxes
func CartFilterByTaxes(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]
		c, _ := r.Cookie("CID")

		minTaxes, _ := strconv.ParseFloat(min, 32)
		maxTaxes, _ := strconv.ParseFloat(max, 32)

		products, err := shopping.FilterByTaxes(db, c.Value, float32(minTaxes), float32(maxTaxes))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByTotal returns the products filtered by total
func CartFilterByTotal(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]
		c, _ := r.Cookie("CID")

		minTotal, _ := strconv.ParseFloat(min, 32)
		maxTotal, _ := strconv.ParseFloat(max, 32)

		products, err := shopping.FilterByTotal(db, c.Value, float32(minTotal), float32(maxTotal))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByType returns the products filtered by type
func CartFilterByType(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productType := mux.Vars(r)["type"]
		c, _ := r.Cookie("CID")

		products, err := shopping.FilterByType(db, c.Value, productType)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByWeight returns the products filtered by weight
func CartFilterByWeight(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]
		c, _ := r.Cookie("CID")

		minWeight, _ := strconv.ParseFloat(min, 32)
		maxWeight, _ := strconv.ParseFloat(max, 32)

		products, err := shopping.FilterByWeight(db, c.Value, float32(minWeight), float32(maxWeight))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartGet returns the cart in a JSON format
func CartGet(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		cart, err := shopping.Get(db, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, cart)
	}
}

// CartItems prints cart items
func CartItems(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")
		items, err := shopping.Items(db, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, items)
	}
}

// CartRemove takes out a product from the shopping cart
func CartRemove(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		q := mux.Vars(r)["quantity"]
		c, _ := r.Cookie("CID")

		key, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = shopping.Remove(db, c.Value, key, quantity)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Successfully removed the product from the cart")
	}
}

// CartReset resets the cart to its default state
func CartReset(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		err := shopping.Reset(db, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Cart reseted")
	}
}

// CartSize returns the size of the shopping cart
func CartSize(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")

		size, err := shopping.Size(db, c.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, size)
	}
}
