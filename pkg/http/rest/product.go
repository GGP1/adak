package rest

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/pkg/product"
	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi"
)

// ProductCreate creates a new product and saves it.
func (s *Frontend) ProductCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p product.Product
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, &p); err != nil {
			http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
			return
		}

		_, err := s.productClient.Create(ctx, &product.CreateRequest{Product: &p})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusCreated, &p)
	}
}

// ProductDelete removes a product.
func (s *Frontend) ProductDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		_, err := s.productClient.Delete(ctx, &product.DeleteRequest{ID: id})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, http.StatusOK, "Product deleted successfully.")
	}
}

// ProductGet lists all the products.
func (s *Frontend) ProductGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		get, err := s.productClient.Get(ctx, &product.GetRequest{})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, get.Products)
	}
}

// ProductGetByID lists the product with the id requested.
func (s *Frontend) ProductGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		getByID, err := s.productClient.GetByID(ctx, &product.GetByIDRequest{ID: id})
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, getByID.Product)
	}
}

// ProductSearch looks for the products with the given value.
func (s *Frontend) ProductSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")
		ctx := r.Context()

		if err := sanitize.Normalize(&query); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		search, err := s.productClient.Search(ctx, &product.SearchRequest{Search: query})
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, search.Products)
	}
}

// ProductUpdate updates the product with the given id.
func (s *Frontend) ProductUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p product.Product
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		_, err := s.productClient.Update(ctx, &product.UpdateRequest{Product: &p, ID: id})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, &p)
	}
}
