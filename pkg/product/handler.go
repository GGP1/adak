package product

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"

	"github.com/go-chi/chi"
	validator "github.com/go-playground/validator/v10"
)

// Handler handles product endpoints.
type Handler struct {
	Service Service
}

// Create creates a new product and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Product
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, product); err != nil {
			http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
			return
		}

		if err := h.Service.Create(ctx, &product); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, product)
	}
}

// Delete removes a product.
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := h.Service.Delete(ctx, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, http.StatusOK, "Product deleted successfully.")
	}
}

// Get lists all the products.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		products, err := h.Service.Get(ctx)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// GetByID lists the product with the id requested.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		product, err := h.Service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, product)
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

		products, err := h.Service.Search(ctx, query)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

// Update updates the product with the given id.
func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Product
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := h.Service.Update(ctx, &product, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, product)
	}
}
