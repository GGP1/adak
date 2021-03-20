package review

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/response"

	"github.com/go-chi/chi"
	validator "github.com/go-playground/validator/v10"
	lru "github.com/hashicorp/golang-lru"
)

// Handler handles reviews endpoints.
type Handler struct {
	Service Service
	Cache   *lru.Cache
}

// Create creates a new review and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review Review
		ctx := r.Context()

		userID, err := cookie.GetValue(r, "UID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, review); err != nil {
			response.Error(w, http.StatusBadRequest, err.(validator.ValidationErrors))
			return
		}

		if err := h.Service.Create(ctx, &review, userID); err != nil {
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

		if err := h.Service.Delete(ctx, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, fmt.Sprintf("review %q deleted", id))
	}
}

// Get lists all the reviews.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reviews, err := h.Service.Get(ctx)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, reviews)
	}
}

// GetByID lists the review with the id requested.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if cReview, ok := h.Cache.Get(id); ok {
			response.JSON(w, http.StatusOK, cReview)
			return
		}

		review, err := h.Service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		h.Cache.Add(id, review)
		response.JSON(w, http.StatusOK, review)
	}
}
