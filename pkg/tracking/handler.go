package tracking

import (
	"net/http"

	"github.com/GGP1/palo/internal/response"

	"github.com/go-chi/chi"
)

// Handler handles tracking endpoints.
type Handler struct {
	TrackerSv Tracker
}

// DeleteHit prints the hit with the specified day.
func (h *Handler) DeleteHit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		if err := h.TrackerSv.Delete(ctx, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Successfully deleted the hit")
	}
}

// GetHits retrieves total amount of hits stored.
func (h *Handler) GetHits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		hits, err := h.TrackerSv.Get(ctx)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, hits)
	}
}

// SearchHit returns the hits that matched with the search.
func (h *Handler) SearchHit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")

		ctx := r.Context()

		hits, err := h.TrackerSv.Search(ctx, query)
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, hits)
	}
}

// SearchHitByField returns the hits that matched with the search.
func (h *Handler) SearchHitByField() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		field := chi.URLParam(r, "field")
		value := chi.URLParam(r, "value")

		ctx := r.Context()

		hits, err := h.TrackerSv.SearchByField(ctx, field, value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, hits)
	}
}
