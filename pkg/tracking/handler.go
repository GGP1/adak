package tracking

import (
	"fmt"
	"net/http"

	"github.com/GGP1/adak/internal/response"

	"github.com/go-chi/chi/v5"
)

// Handler handles tracking endpoints.
type Handler struct {
	service Tracker
}

// NewHandler returns a new tracking handler.
func NewHandler(service Tracker) Handler {
	return Handler{service}
}

// DeleteHit prints the hit with the specified day.
func (h *Handler) DeleteHit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := h.service.Delete(ctx, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, fmt.Sprintf("hit %q deleted", id))
	}
}

// GetHits retrieves total amount of hits stored.
func (h *Handler) GetHits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		hits, err := h.service.Get(ctx)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, hits)
	}
}

// SearchHit returns the hits that matched with the search.
func (h *Handler) SearchHit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")
		ctx := r.Context()

		hits, err := h.service.Search(ctx, query)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, hits)
	}
}

// SearchHitByField returns the hits that matched with the search.
func (h *Handler) SearchHitByField() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		field := chi.URLParam(r, "field")
		value := chi.URLParam(r, "value")
		ctx := r.Context()

		hits, err := h.service.SearchByField(ctx, field, value)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, hits)
	}
}
