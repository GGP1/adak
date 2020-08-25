package review

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/internal/token"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// Reviews handles reviews routes.
type Reviews struct {
	DB *sqlx.DB
}

// Create creates a new review and saves it.
func Create(rev Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID, _ := r.Cookie("UID")

		var (
			review Review
			ctx    = r.Context()
		)

		userID, err := token.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := rev.Create(ctx, &review, userID); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusCreated, review)
	}
}

// Delete removes a review.
func Delete(rev Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		if err := rev.Delete(ctx, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Review deleted successfully.")
	}
}

// Get lists all the reviews.
func Get(rev Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reviews, err := rev.Get(ctx)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, reviews)
	}
}

// GetByID lists the review with the id requested.
func GetByID(rev Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		ctx := r.Context()

		review, err := rev.GetByID(ctx, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}
