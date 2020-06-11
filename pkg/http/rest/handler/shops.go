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
	"github.com/jinzhu/gorm"
)

// ShopHandler defines all of the handlers related to shops. It holds the
// application state needed by the handler methods.
type ShopHandler struct {
	DB *gorm.DB
}

// Get lists all the shops
func (sh *ShopHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop []model.Shop

		err := listing.GetShops(&shop, sh.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, shop)
	}
}

// GetOne lists the shop with the id requested
func (sh *ShopHandler) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetAShop(&shop, id, sh.DB)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Shop not found")
			return
		}

		response.RespondJSON(w, r, http.StatusOK, shop)
	}
}

// Add creates a new shop and saves it
func (sh *ShopHandler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := adding.AddShop(&shop, sh.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, shop)
	}
}

// Update updates the shop with the given id
func (sh *ShopHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		param := mux.Vars(r)
		id := param["id"]

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := updating.UpdateShop(&shop, id, sh.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Shop updated")
	}
}

// Delete deletes a shop
func (sh *ShopHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteShop(&shop, id, sh.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Shop deleted")
	}
}
