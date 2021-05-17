package product

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/internal/validate"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
)

// Handler handles product endpoints.
type Handler struct {
	Service Service
	Cache   *memcache.Client
}

// Create creates a new product and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var p Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validate.Struct(ctx, p); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		p.ID = zero.StringFrom(token.RandString(27))
		p.CreatedAt = zero.TimeFrom(time.Now())
		if err := h.Service.Create(ctx, p); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, p)
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

		response.JSONText(w, http.StatusOK, id)
	}
}

// Get lists all the products.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		products, err := h.Service.Get(ctx)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
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

		item, err := h.Cache.Get(id)
		if err == nil {
			response.EncodedJSON(w, item.Value)
			return
		}

		product, err := h.Service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSONAndCache(h.Cache, w, id, product)
	}
}

// Search looks for the products with the given value.
func (h *Handler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")
		ctx := r.Context()

		query = sanitize.Normalize(query)
		if strings.ContainsAny(query, ";-\\|@#~€¬<>_()[]}{¡^'") {
			response.Error(w, http.StatusBadRequest, errors.Errorf("query contains invalid characters"))
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
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		var product UpdateProduct
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := h.Service.Update(ctx, id, product); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, product)
	}
}
