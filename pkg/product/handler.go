package product

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"

	"github.com/go-chi/chi"
)

// Handler handles product endpoints.
type Handler struct {
	Service Service
}

// Create creates a new product and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			product Product
			ctx     = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := h.Service.Create(ctx, &product); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusCreated, product)
	}
}

// Delete removes a product.
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		if err := h.Service.Delete(ctx, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Product deleted successfully.")
	}
}

// Get lists all the products.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		products, err := h.Service.Get(ctx)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// GetByID lists the product with the id requested.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		product, err := h.Service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}

// Search looks for the products with the given value.
func (h *Handler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")

		ctx := r.Context()

		products, err := h.Service.Search(ctx, query)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// Update updates the product with the given id.
func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var (
			product Product
			ctx     = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := h.Service.Update(ctx, &product, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}
