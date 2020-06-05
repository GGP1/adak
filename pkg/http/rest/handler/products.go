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

	"github.com/gorilla/mux"
)

// GetProducts lists all the products
func GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product []model.Product

		err := listing.GetProducts(&product)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, product)
	}
}

// GetOneProduct lists one product based on the id
func GetOneProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		param := mux.Vars(r)
		id := param["id"]

		err := listing.GetAProduct(&product, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Product not found")
			return
		}

		response.RespondJSON(w, r, http.StatusOK, product)
	}
}

// AddProduct creates a new product and saves it
func AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}

		err := adding.AddProduct(&product)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusCreated, product)
	}
}

// UpdateProduct updates a product
func UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product

		param := mux.Vars(r)
		id := param["id"]

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Respond(w, r, http.StatusInternalServerError, err)
		}
		defer r.Body.Close()

		err := updating.UpdateProduct(&product, id)
		if err != nil {
			response.Respond(w, r, http.StatusNotFound, err)
		}

		response.RespondJSON(w, r, http.StatusOK, product)
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
			response.Respond(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Product deleted")
	}
}
