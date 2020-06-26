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

	"github.com/gorilla/mux"
)

// GetProducts lists all the products
func GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product []model.Product

		err := listing.GetProducts(&product)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}

// GetOneProduct lists the product with the id requested
func GetOneProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetOneProduct(&product, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}

// AddProduct creates a new product and saves it
func AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := adding.AddProduct(&product)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusCreated, product)
	}
}

// UpdateProduct updates the product with the given id
func UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		param := mux.Vars(r)
		id := param["id"]

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		err := updating.UpdateProduct(&product, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusOK, product)
	}
}

// DeleteProduct deletes a product
func DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		param := mux.Vars(r)
		id := param["id"]

		err := deleting.DeleteProduct(&product, id)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Product deleted successfully.")
	}
}
