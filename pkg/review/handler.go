package review

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/params"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/internal/validate"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
	"gopkg.in/guregu/null.v4/zero"
)

type cursorResponse struct {
	NextCursor string   `json:"next_cursor,omitempty"`
	Reviews    []Review `json:"reviews,omitempty"`
}

// Handler handles reviews endpoints.
type Handler struct {
	service Service
	cache   *memcache.Client
}

// NewHandler returns a new review handler.
func NewHandler(service Service, cache *memcache.Client) Handler {
	return Handler{
		service: service,
		cache:   cache,
	}
}

// Create creates a new review and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := cookie.GetValue(r, "UID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		if err := token.CheckPermits(r, userID); err != nil {
			response.Error(w, http.StatusUnauthorized, err)
			return
		}

		var review Review
		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validate.Struct(ctx, review); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		review.ID = zero.StringFrom(token.RandString(28))
		if err := h.service.Create(ctx, review); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, review)
	}
}

// Delete removes a review.
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

// Get lists all the reviews.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		urlParams, err := params.ParseQuery(r.URL.RawQuery, params.Review)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		reviews, err := h.service.Get(ctx, urlParams)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		var nextCursor string
		if len(reviews) > 0 {
			nextCursor = params.EncodeCursor(
				reviews[len(reviews)-1].CreatedAt.Time,
				reviews[len(reviews)-1].ID.String,
			)
		}

		response.JSON(w, http.StatusOK, cursorResponse{
			NextCursor: nextCursor,
			Reviews:    reviews,
		})
	}
}

// GetByID lists the review with the id requested.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		item, err := h.cache.Get(id)
		if err == nil {
			response.EncodedJSON(w, item.Value)
			return
		}

		review, err := h.service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSONAndCache(h.cache, w, id, review)
	}
}
