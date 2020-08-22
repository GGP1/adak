package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/creating"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
)

// Reviews handles reviews routes.
type Reviews struct {
	DB *sqlx.DB
}

// CreateReview creates a new review and saves it.
func CreateReview(c creating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		uID, _ := r.Cookie("UID")

		userID, err := auth.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := c.CreateReview(&review, userID); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusCreated, review)
	}
}

// DeleteReview removes a review.
func DeleteReview(d deleting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := d.DeleteReview(id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Review deleted successfully.")
	}
}

// GetReviews lists all the reviews.
func GetReviews(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reviews, err := l.GetReviews()
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, reviews)
	}
}

// GetReviewByID lists the review with the id requested.
func GetReviewByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		review, err := l.GetReviewByID(id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}
