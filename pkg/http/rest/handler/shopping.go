package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/gorilla/mux"
)

// CartAdd appends a product to the cart
func CartAdd(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

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

		err = cart.Add(&product, quantity)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		response.JSON(w, r, http.StatusOK, cart)
	}
}

// CartCheckout returns the final purchase
func CartCheckout(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkout := cart.Checkout()

		response.JSON(w, r, http.StatusOK, checkout)
	}
}

// CartFilterByBrand returns the products filtered by brand
func CartFilterByBrand(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		brand := mux.Vars(r)["brand"]

		products, err := cart.FilterByBrand(brand)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByCategory returns the products filtered by category
func CartFilterByCategory(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := mux.Vars(r)["category"]

		products, err := cart.FilterByCategory(category)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByDiscount returns the products filtered by discount
func CartFilterByDiscount(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]

		minDiscount, _ := strconv.ParseFloat(min, 32)
		maxDiscount, _ := strconv.ParseFloat(max, 32)

		products, err := cart.FilterByDiscount(float32(minDiscount), float32(maxDiscount))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterBySubtotal returns the products filtered by subtotal
func CartFilterBySubtotal(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]

		minSubtotal, _ := strconv.ParseFloat(min, 32)
		maxSubtotal, _ := strconv.ParseFloat(max, 32)

		products, err := cart.FilterBySubtotal(float32(minSubtotal), float32(maxSubtotal))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByTaxes returns the products filtered by taxes
func CartFilterByTaxes(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]

		minTaxes, _ := strconv.ParseFloat(min, 32)
		maxTaxes, _ := strconv.ParseFloat(max, 32)

		products, err := cart.FilterByTaxes(float32(minTaxes), float32(maxTaxes))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByTotal returns the products filtered by total
func CartFilterByTotal(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]

		minTotal, _ := strconv.ParseFloat(min, 32)
		maxTotal, _ := strconv.ParseFloat(max, 32)

		products, err := cart.FilterByTotal(float32(minTotal), float32(maxTotal))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByType returns the products filtered by type
func CartFilterByType(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productType := mux.Vars(r)["type"]

		products, err := cart.FilterByType(productType)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByWeight returns the products filtered by weight
func CartFilterByWeight(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := mux.Vars(r)["min"]
		max := mux.Vars(r)["max"]

		minWeight, _ := strconv.ParseFloat(min, 32)
		maxWeight, _ := strconv.ParseFloat(max, 32)

		products, err := cart.FilterByWeight(float32(minWeight), float32(maxWeight))
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartGet returns the cart in a JSON format
func CartGet(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, r, http.StatusOK, cart)
	}
}

// CartItems prints cart items
func CartItems(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := cart.Items()

		response.JSON(w, r, http.StatusOK, items)
	}
}

// CartRemove takes out a product from the shopping cart
func CartRemove(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		q := mux.Vars(r)["quantity"]

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

		err = cart.Remove(uint(key), quantity)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		response.JSON(w, r, http.StatusOK, cart)
	}
}

// CartReset resets the cart to its default state
func CartReset(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cart.Reset()

		response.JSON(w, r, http.StatusOK, cart)
	}
}

// CartSize returns the size of the shopping cart
func CartSize(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		size := cart.Size()

		response.JSON(w, r, http.StatusOK, size)
	}
}
