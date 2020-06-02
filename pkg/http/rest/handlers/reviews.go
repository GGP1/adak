package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/GGP1/palo/internal/utils/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/models"
	"github.com/gorilla/mux"
)

// GetReviews lists all the reviews
func GetReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review []models.Review

		err := listing.GetReviews(&review)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.Respond(w, r, http.StatusOK, review)
	}
}

// GetSingleReview lists a review based on the id
func GetOneReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review models.Review

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetAReview(&review, id)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		if review.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Review not found")
			return
		}

		response.Respond(w, r, http.StatusOK, review)
	}
}

// AddReview creates a new review and saves it
func AddReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review models.Review
		var buf bytes.Buffer
		var err error

		err = json.NewEncoder(&buf).Encode(&review)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Review not found")
		}

		err = adding.AddReview(&review)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Review deleted")
	}
}

// DeleteReview deletes a review
func DeleteReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review models.Review

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteReview(&review, id)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.Respond(w, r, http.StatusOK, review)
	}
}
