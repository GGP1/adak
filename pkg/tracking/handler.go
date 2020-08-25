package tracking

import (
	"net/http"

	"github.com/GGP1/palo/internal/response"

	"github.com/go-chi/chi"
)

// DeleteHit prints the hit with the specified day.
func DeleteHit(t Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		if err := t.Delete(ctx, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Successfully deleted the hit")
	}
}

// GetHits retrieves total amount of hits stored.
func GetHits(t Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		hits, err := t.Get(ctx)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, hits)
	}
}

// SearchHit returns the hits that matched with the search.
func SearchHit(t Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := chi.URLParam(r, "search")

		ctx := r.Context()

		hits, err := t.Search(ctx, search)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, hits)
	}
}

// SearchHitByField returns the hits that matched with the search.
func SearchHitByField(t Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		field := chi.URLParam(r, "field")
		value := chi.URLParam(r, "value")

		ctx := r.Context()

		hits, err := t.SearchByField(ctx, field, value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, hits)
	}
}
