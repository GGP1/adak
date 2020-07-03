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
)

// GetReviews lists all the reviews
func GetReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review []model.Review

		err := listing.GetReviews(&review)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}

// GetReviewByID lists the review with the id requested
func GetReviewByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetReviewByID(&review, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}

// AddReview creates a new review and saves it
func AddReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := adding.AddReview(&review)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review)
	}
}

// DeleteReview deletes a review
func DeleteReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteReview(&review, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Review deleted successfully.")
	}
}
