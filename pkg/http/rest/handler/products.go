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
	"github.com/GGP1/palo/pkg/updating"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

// ProductHandler defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type ProductHandler struct {
	DB *gorm.DB
}

// Get lists all the products
func (ph *ProductHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product []model.Product

		err := listing.GetProducts(&product, ph.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, product)
	}
}

// GetOne lists one product based on the id
func (ph *ProductHandler) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetAProduct(&product, id, ph.DB)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Product not found")
			return
		}

		response.RespondJSON(w, r, http.StatusOK, product)
	}
}

// Add creates a new product and saves it
func (ph *ProductHandler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := adding.AddProduct(&product, ph.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusCreated, product)
	}
}

// Update updates a product
func (ph *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		param := mux.Vars(r)
		id := param["id"]

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := updating.UpdateProduct(&product, id, ph.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, product)
	}
}

// Delete deletes a product
func (ph *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteProduct(&product, id, ph.DB)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Product deleted")
	}
}
