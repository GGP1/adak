package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/updating"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

// Products handles products routes.
type Products struct {
	DB *gorm.DB
}

// Add creates a new product and saves it.
func (p *Products) Add(a adding.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := a.AddProduct(p.DB, &product)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}

// Delete removes a product.
func (p *Products) Delete(d deleting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		id := mux.Vars(r)["id"]

		err := d.DeleteProduct(p.DB, &product, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Product deleted successfully.")
	}
}

// Find lists all the products.
func (p *Products) Find(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product []model.Product

		err := l.GetProducts(p.DB, &product)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}

// FindByID lists the product with the id requested.
func (p *Products) FindByID(l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		id := mux.Vars(r)["id"]

		err := l.GetProductByID(p.DB, &product, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}

// Update updates the product with the given id.
func (p *Products) Update(u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		id := mux.Vars(r)["id"]

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := u.UpdateProduct(p.DB, &product, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}
