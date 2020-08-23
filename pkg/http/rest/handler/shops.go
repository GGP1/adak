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

	"github.com/go-chi/chi"
)

// CreateShop creates a new shop and saves it.
func CreateShop(c creating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			shop *model.Shop
			ctx  = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := c.CreateShop(ctx, shop); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusCreated, shop)
	}
}

// DeleteShop removes a shop.
func DeleteShop(d deleting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		if err := d.DeleteShop(ctx, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Shop deleted successfully.")
	}
}

// GetShops lists all the shops.
func GetShops(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		shops, err := l.GetShops(ctx)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shops)
	}
}

// GetShopByID lists the shop with the id requested.
func GetShopByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		shop, err := l.GetShopByID(ctx, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shop)
	}
}

// SearchShop looks for the products with the given value.
func SearchShop(s searching.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")

		ctx := r.Context()

		shops, err := s.SearchShops(ctx, query)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, shops)
	}
}

// UpdateShop updates the shop with the given id.
func UpdateShop(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var (
			shop *model.Shop
			ctx  = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := u.UpdateShop(ctx, shop, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Shop updated successfully.")
	}
}
