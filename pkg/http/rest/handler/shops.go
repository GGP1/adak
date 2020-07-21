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
	"github.com/jinzhu/gorm"
)

// Shops handles shops routes
type Shops struct {
	DB *gorm.DB
}

// Add creates a new shop and saves it
func (s *Shops) Add(a adding.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := a.AddShop(s.DB, &shop)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shop)
	}
}

// Delete deletes a shop
func (s *Shops) Delete(d deleting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		id := mux.Vars(r)["id"]

		err := d.DeleteShop(s.DB, &shop, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Shop deleted successfully.")
	}
}

// GetAll lists all the shops
func (s *Shops) GetAll(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop []model.Shop

		err := l.GetShops(s.DB, &shop)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shop)
	}
}

// GetByID lists the shop with the id requested
func (s *Shops) GetByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		id := mux.Vars(r)["id"]

		err := l.GetShopByID(s.DB, &shop, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shop)
	}
}

// Update updates the shop with the given id
func (s *Shops) Update(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		id := mux.Vars(r)["id"]

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := u.UpdateShop(s.DB, &shop, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Shop updated successfully.")
	}
}
