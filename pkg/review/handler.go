package review

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/token"

	"github.com/go-chi/chi"
	validator "github.com/go-playground/validator/v10"
)

// Handler handles reviews endpoints.
type Handler struct {
	Service Service
}

// Create creates a new review and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review Review
		ctx := r.Context()
		uID, err := r.Cookie("UID")
		if err != nil {
			http.Redirect(w, r, "http://localhost:4000/login", http.StatusTemporaryRedirect)
			return
		}

		userID, err := token.GetUserID(uID.Value)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, review); err != nil {
			http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
			return
		}

		if err := h.Service.Create(ctx, &review, userID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, review)
	}
}

// Delete removes a review.
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := h.Service.Delete(ctx, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, http.StatusOK, "Review deleted successfully.")
	}
}

// Get lists all the reviews.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reviews, err := h.Service.Get(ctx)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, reviews)
	}
}

// GetByID lists the review with the id requested.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		review, err := h.Service.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, review)
	}
}
