package shop

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"

	"github.com/go-chi/chi"
	validator "github.com/go-playground/validator/v10"
)

// Handler handles shop endpoints.
type Handler struct {
	Service Service
}

// Create creates a new shop and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop *Shop
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, shop); err != nil {
			http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
			return
		}

		if err := h.Service.Create(ctx, shop); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, shop)
	}
}

// Delete removes a shop.
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := h.Service.Delete(ctx, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, http.StatusOK, "Shop deleted successfully.")
	}
}

// Get lists all the shops.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		shops, err := h.Service.Get(ctx)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, shops)
	}
}

// GetByID lists the shop with the id requested.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		shop, err := h.Service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, shop)
	}
}

// Search looks for the products with the given value.
func (h *Handler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")
		ctx := r.Context()

		if err := sanitize.Normalize(&query); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		shops, err := h.Service.Search(ctx, query)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, shops)
	}
}

// Update updates the shop with the given id.
func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shop *Shop
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := h.Service.Update(ctx, shop, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, http.StatusOK, "Shop updated successfully.")
	}
}
