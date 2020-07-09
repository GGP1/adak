package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/updating"
	"github.com/gorilla/mux"
)

// GetShops lists all the shops
func GetShops(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop []model.Shop

		err := l.GetShops(&shop)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shop)
	}
}

// GetShopByID lists the shop with the id requested
func GetShopByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		id := mux.Vars(r)["id"]

		err := l.GetShopByID(&shop, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shop)
	}
}

// AddShop creates a new shop and saves it
func AddShop(a adding.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := a.AddShop(&shop)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shop)
	}
}

// UpdateShop updates the shop with the given id
func UpdateShop(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		id := mux.Vars(r)["id"]

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := u.UpdateShop(&shop, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Shop updated successfully.")
	}
}

// DeleteShop deletes a shop
func DeleteShop(d deleting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		id := mux.Vars(r)["id"]

		err := d.DeleteShop(&shop, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		response.HTMLText(w, r, http.StatusOK, "Shop deleted successfully.")
	}
}
