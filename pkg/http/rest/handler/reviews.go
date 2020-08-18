package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/creating"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

// Reviews handles reviews routes.
type Reviews struct {
	DB *gorm.DB
}

// Create creates a new review and saves it.
func (rev *Reviews) Create(c creating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := c.CreateReview(rev.DB, &review); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}

// Delete removes a review.
func (rev *Reviews) Delete(d deleting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := d.DeleteReview(rev.DB, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Review deleted successfully.")
	}
}

// Get lists all the reviews.
func (rev *Reviews) Get(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review []model.Review

		if err := l.GetReviews(rev.DB, &review); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}

// GetByID lists the review with the id requested.
func (rev *Reviews) GetByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var review model.Review

		if err := l.GetReviewByID(rev.DB, &review, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}
