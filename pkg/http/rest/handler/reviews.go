package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GGP1/palo/internal/utils/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// ReviewHandler defines all of the handlers related to reviews. It holds the
// application state needed by the handler methods.
type ReviewHandler struct {
	DB *gorm.DB
}

// Get lists all the reviews
func (rh *ReviewHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review []model.Review

		err := listing.GetReviews(&review, rh.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, review)
	}
}

// GetOne lists the review with the id requested
func (rh *ReviewHandler) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetAReview(&review, id, rh.DB)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Review not found")
			return
		}

		response.RespondJSON(w, r, http.StatusOK, review)
	}
}

// Add creates a new review and saves it
func (rh *ReviewHandler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := adding.AddReview(&review, rh.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Review deleted")
	}
}

// Delete deletes a review
func (rh *ReviewHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review model.Review

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteReview(&review, id, rh.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Review deleted")
	}
}
