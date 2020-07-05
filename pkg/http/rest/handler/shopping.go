package handler

import (
	"encoding/json"
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

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		cart.Add(&product)
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

// CartRemove takes out a product from the shopping cart
func CartRemove(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		key, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		cart.Remove(uint(key))

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

// CartShowItems prints cart items
func CartShowItems(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := cart.ShowItems()

		response.JSON(w, r, http.StatusOK, items)
	}
}

// CartSize returns the size of the shopping cart
func CartSize(cart *shopping.Cart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		size := cart.Size()

		response.JSON(w, r, http.StatusOK, size)
	}
}
