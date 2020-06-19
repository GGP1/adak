package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GGP1/palo/internal/utils/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/updating"
	"github.com/gorilla/mux"
)

// GetShops lists all the shops
func GetShops() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop []model.Shop

		err := listing.GetAll(&shop)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, shop)
	}
}

// GetOneShop lists the shop with the id requested
func GetOneShop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetOne(&shop, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Shop not found")
			return
		}

		response.RespondJSON(w, r, http.StatusOK, shop)
	}
}

// AddShop creates a new shop and saves it
func AddShop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := adding.Add(&shop)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, shop)
	}
}

// UpdateShop updates the shop with the given id
func UpdateShop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		param := mux.Vars(r)
		id := param["id"]

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

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
		var shop model.Shop

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.Delete(&shop, id)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Shop deleted")
	}
}
