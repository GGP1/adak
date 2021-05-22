package shop

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/internal/validate"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// Handler handles shop endpoints.
type Handler struct {
	service Service
	cache   *memcache.Client
}

// NewHandler returns a new shop handler.
func NewHandler(service Service, cache *memcache.Client) Handler {
	return Handler{
		service: service,
		cache:   cache,
	}
}

// Create creates a new shop and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var shop Shop
		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validate.Struct(ctx, shop); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		shop.ID = token.RandString(29)
		if err := h.service.Create(ctx, shop); err != nil {
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

		if err := h.service.Delete(ctx, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, id)
	}
}

// Get lists all the shops.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		shops, err := h.service.Get(ctx)
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

		item, err := h.cache.Get(id)
		if err == nil {
			response.EncodedJSON(w, item.Value)
			return
		}

		shop, err := h.service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSONAndCache(h.cache, w, id, shop)
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

		shops, err := h.service.Search(ctx, query)
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
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		var shop UpdateShop
		if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := h.service.Update(ctx, id, shop); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, id)
	}
}
