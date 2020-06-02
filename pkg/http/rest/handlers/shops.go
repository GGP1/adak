package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/GGP1/palo/internal/utils/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/models"
	"github.com/GGP1/palo/pkg/updating"
	"github.com/gorilla/mux"
)

// GetShops lists all the shops
func GetShops() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop []models.Shop

		err := listing.GetShops(&shop)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.Respond(w, r, http.StatusOK, shop)
	}
}

// GetOneShop lists one shop based on the id
func GetOneShop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop models.Shop

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetAShop(&shop, id)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		if shop.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Shop not found")
			return
		}

		response.Respond(w, r, http.StatusOK, shop)
	}
}

// AddShop creates a new shop and saves it
func AddShop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop models.Shop
		var buf bytes.Buffer
		var err error

		err = json.NewEncoder(&buf).Encode(&shop)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Review not found")
		}

		err = adding.AddShop(&shop)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.Respond(w, r, http.StatusOK, shop)
	}
}

// UpdateShop updates a shop
func UpdateShop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop models.Shop

		param := mux.Vars(r)
		id := param["id"]

		err := updating.UpdateShop(&shop, id)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Shop updated")
	}
}

// DeleteShop deletes a shop
func DeleteShop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop models.Shop

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteShop(&shop, id)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Shop deleted")
	}
}
