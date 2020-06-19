package handler

import (
	"encoding/json"
	"io"
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

		err := listing.GetAll(&review)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, review)
	}
}

// GetOneReview lists the review with the id requested
func GetOneReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetOne(&review, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Review not found")
			return
		}

		response.RespondJSON(w, r, http.StatusOK, review)
	}
}

// AddReview creates a new review and saves it
func AddReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := adding.Add(&review)
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
		var review model.Review

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.Delete(&review, id)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Review deleted")
	}
}
