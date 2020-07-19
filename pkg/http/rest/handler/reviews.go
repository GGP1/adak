package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Reviews handles reviews routes
type Reviews struct {
	DB *gorm.DB
}

// Add creates a new review and saves it
func (rev *Reviews) Add(a adding.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := a.AddReview(rev.DB, &review)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}

// Delete deletes a review
func (rev *Reviews) Delete(d deleting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		id := mux.Vars(r)["id"]

		err := d.DeleteReview(rev.DB, &review, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Review deleted successfully.")
	}
}

// GetAll lists all the reviews
func (rev *Reviews) GetAll(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review []model.Review

		err := l.GetReviews(rev.DB, &review)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}

// GetByID lists the review with the id requested
func (rev *Reviews) GetByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		id := mux.Vars(r)["id"]

		err := l.GetReviewByID(rev.DB, &review, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}
