package product

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/params"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/validate"
	"github.com/google/uuid"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4/zero"
)

type cursorResponse struct {
	NextCursor string    `json:"next_cursor,omitempty"`
	Products   []Product `json:"products,omitempty"`
}

// Handler handles product endpoints.
type Handler struct {
	service Service
	cache   *memcache.Client
}

// NewHandler returns a new product handler.
func NewHandler(service Service, cache *memcache.Client) Handler {
	return Handler{
		service: service,
		cache:   cache,
	}
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

		p.ID = zero.StringFrom(uuid.NewString())
		p.CreatedAt = zero.TimeFrom(time.Now())
		if err := h.service.Create(ctx, p); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, p)
	}
}

// Delete removes a product.
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := params.URLID(ctx)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		if err := h.service.Delete(ctx, id); err != nil {
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

		urlParams, err := params.ParseQuery(r.URL.RawQuery, params.Product)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		products, err := h.service.Get(ctx, urlParams)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		var nextCursor string
		if len(products) > 0 {
			nextCursor = params.EncodeCursor(
				products[len(products)-1].CreatedAt.Time,
				products[len(products)-1].ID.String,
			)
		}

		response.JSON(w, http.StatusOK, cursorResponse{
			NextCursor: nextCursor,
			Products:   products,
		})
	}
}

// GetByID lists the product with the id requested.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := params.URLID(ctx)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		item, err := h.cache.Get(id)
		if err == nil {
			response.EncodedJSON(w, item.Value)
			return
		}

		product, err := h.service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSONAndCache(h.cache, w, id, product)
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

		products, err := h.service.Search(ctx, query)
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
		ctx := r.Context()

		id, err := params.URLID(ctx)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		var product UpdateProduct
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validate.Struct(ctx, product); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		if err := h.service.Update(ctx, id, product); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, product)
	}
}
