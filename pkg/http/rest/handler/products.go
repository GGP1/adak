package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/creating"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/searching"
	"github.com/GGP1/palo/pkg/updating"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// Products handles products routes.
type Products struct {
	DB *sqlx.DB
}

// Create creates a new product and saves it.
func (p *Products) Create(c creating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := c.CreateProduct(p.DB, &product); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}

// Delete removes a product.
func (p *Products) Delete(d deleting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if err := d.DeleteProduct(p.DB, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Product deleted successfully.")
	}
}

// Get lists all the products.
func (p *Products) Get(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := l.GetProducts(p.DB)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// GetByID lists the product with the id requested.
func (p *Products) GetByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		product, err := l.GetProductByID(p.DB, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}

// Search looks for the products with the given value.
func (p *Products) Search(s searching.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := chi.URLParam(r, "search")

		var products []model.Product

		if err := s.SearchProducts(p.DB, &products, search); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// Update updates the product with the given id.
func (p *Products) Update(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var product model.Product

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := u.UpdateProduct(p.DB, &product, id); err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}
