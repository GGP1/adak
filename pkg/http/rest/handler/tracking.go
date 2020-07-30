package handler

import (
	"context"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/tracking"
	"github.com/gorilla/mux"
)

// DeleteHit prints the hit with the specified day.
func DeleteHit(ctx context.Context, t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		err := t.DeleteHit(ctx, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Successfully deleted the hit")
	}
}

// GetHits retrieves total amount of hits stored.
func GetHits(ctx context.Context, t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hits, err := t.Get(ctx)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, hits)
	}
}

// SearchHit returns the hits that matched with the search
func SearchHit(ctx context.Context, t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := mux.Vars(r)["search"]

		hits, err := t.Search(ctx, search)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, hits)
	}
}
