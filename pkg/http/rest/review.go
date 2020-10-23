package rest

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/internal/sanitize"
	"github.com/GGP1/palo/internal/token"
	"github.com/GGP1/palo/pkg/review"
	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi"
)

// ReviewCreate creates a new review and saves it.
func (s *Frontend) ReviewCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID, _ := r.Cookie("UID")
		ctx := r.Context()

		var rw review.Review

		userID, err := token.ParseFixedJWT(uID.Value)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&rw); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, &rw); err != nil {
			http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
			return
		}

		if err := sanitize.Normalize(&rw.Comment); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		_, err = s.reviewClient.Create(ctx, &review.CreateRequest{Review: &rw, UserID: userID})
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusCreated, &rw)
	}
}

// ReviewDelete removes a review.
func (s *Frontend) ReviewDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		_, err := s.reviewClient.Delete(ctx, &review.DeleteRequest{ID: id})
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Review deleted successfully.")
	}
}

// ReviewGet lists all the reviews.
func (s *Frontend) ReviewGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reviews, err := s.reviewClient.Get(ctx, &review.GetRequest{})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, reviews.Reviews)
	}
}

// ReviewGetByID lists the review with the id requested.
func (s *Frontend) ReviewGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		review, err := s.reviewClient.GetByID(ctx, &review.GetByIDRequest{ID: id})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, review.Reviews)
	}
}
