package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/creating"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/searching"
	"github.com/GGP1/palo/pkg/updating"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
)

// Shops handles shops routes
type Shops struct {
	DB *sqlx.DB
}

// Create creates a new shop and saves it.
func (s *Shops) Create(c creating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop model.Shop

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := c.CreateShop(s.DB, &shop); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shop)
	}
}

// Delete removes a shop.
func (s *Shops) Delete(d deleting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := d.DeleteShop(s.DB, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Shop deleted successfully.")
	}
}

// Get lists all the shops.
func (s *Shops) Get(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shops, err := l.GetShops(s.DB)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shops)
	}
}

// GetByID lists the shop with the id requested.
func (s *Shops) GetByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		shop, err := l.GetShopByID(s.DB, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shop)
	}
}

// Search looks for the products with the given value.
func (s *Shops) Search(sr searching.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := chi.URLParam(r, "search")

		var shops []model.Shop

		if err := sr.SearchShops(s.DB, &shops, search); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shops)
	}
}

// Update updates the shop with the given id.
func (s *Shops) Update(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var shop model.Shop

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := u.UpdateShop(s.DB, &shop, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Shop updated successfully.")
	}
}
