package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/internal/sanitize"
	"github.com/GGP1/palo/pkg/shopping/cart"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var (
	errInvalidMinNumber = errors.New("invalid minimu number")
	errInvalidMaxNumber = errors.New("invalid max number")
)

// CartAdd appends a product to the cart.
func (s *Frontend) CartAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product *cart.Product
		q := chi.URLParam(r, "quantity")
		ctx := r.Context()

		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
		}

		if quantity == 0 {
			response.Error(w, r, http.StatusBadRequest, errors.New("quantity must be higher than zero"))
			return
		}

		if err = json.NewDecoder(r.Body).Decode(&product); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		cartID, _ := r.Cookie("CID")
		cart, err := s.shoppingClient.Add(ctx, &cart.AddRequest{
			CartID:   cartID.Value,
			Product:  product,
			Quantity: int64(quantity),
		})
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, r, http.StatusCreated, cart)
	}
}

// CartCheckout returns the final purchase.
func (s *Frontend) CartCheckout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		checkout, err := s.shoppingClient.Checkout(ctx, &cart.CheckoutRequest{
			CartID: c.Value,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, checkout)
	}
}

// CartFilterByBrand returns the products filtered by brand.
func (s *Frontend) CartFilterByBrand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		brand := chi.URLParam(r, "brand")
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		if err := sanitize.Normalize(&brand); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		products, err := s.shoppingClient.FilterByBrand(ctx, &cart.FilterTextRequest{
			CartID: c.Value,
			Field:  brand,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByCategory returns the products filtered by category.
func (s *Frontend) CartFilterByCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := chi.URLParam(r, "category")
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		if err := sanitize.Normalize(&category); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		products, err := s.shoppingClient.FilterByCategory(ctx, &cart.FilterTextRequest{
			CartID: c.Value,
			Field:  category,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByDiscount returns the products filtered by discount.
func (s *Frontend) CartFilterByDiscount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		minD, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMinNumber)
			return
		}

		maxD, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMaxNumber)
			return
		}

		products, err := s.shoppingClient.FilterByDiscount(ctx, &cart.FilterNumberRequest{
			CartID: c.Value,
			Min:    minD,
			Max:    maxD,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterBySubtotal returns the products filtered by subtotal.
func (s *Frontend) CartFilterBySubtotal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		minS, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMinNumber)
			return
		}

		maxS, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMaxNumber)
			return
		}

		products, err := s.shoppingClient.FilterBySubtotal(ctx, &cart.FilterNumberRequest{
			CartID: c.Value,
			Min:    minS,
			Max:    maxS,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByTaxes returns the products filtered by taxes.shoppingClient.
func (s *Frontend) CartFilterByTaxes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		minT, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMinNumber)
			return
		}

		maxT, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMaxNumber)
			return
		}

		products, err := s.shoppingClient.FilterByTaxes(ctx, &cart.FilterNumberRequest{
			CartID: c.Value,
			Min:    minT,
			Max:    maxT,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByTotal returns the products filtered by total.
func (s *Frontend) CartFilterByTotal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		minT, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMinNumber)
			return
		}

		maxT, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMaxNumber)
			return
		}

		products, err := s.shoppingClient.FilterByTotal(ctx, &cart.FilterNumberRequest{
			CartID: c.Value,
			Min:    minT,
			Max:    maxT,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByType returns the products filtered by type.
func (s *Frontend) CartFilterByType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pType := chi.URLParam(r, "type")
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		if err := sanitize.Normalize(&pType); err != nil {
			response.Error(w, r, http.StatusBadRequest, err)
			return
		}

		products, err := s.shoppingClient.FilterByType(ctx, &cart.FilterTextRequest{
			CartID: c.Value,
			Field:  pType,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartFilterByWeight returns the products filtered by weight.
func (s *Frontend) CartFilterByWeight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		min := chi.URLParam(r, "min")
		max := chi.URLParam(r, "max")
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		minW, err := strconv.ParseFloat(min, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMinNumber)
			return
		}

		maxW, err := strconv.ParseFloat(max, 64)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errInvalidMaxNumber)
			return
		}

		products, err := s.shoppingClient.FilterByWeight(ctx, &cart.FilterNumberRequest{
			CartID: c.Value,
			Min:    minW,
			Max:    maxW,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, products)
	}
}

// CartGet returns the cart in a JSON format.
func (s *Frontend) CartGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		cart, err := s.shoppingClient.Get(ctx, &cart.GetRequest{CartID: c.Value})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, cart)
	}
}

// CartProducts retrieves cart products.shoppingClient.
func (s *Frontend) CartProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		items, err := s.shoppingClient.Products(ctx, &cart.ProductsRequest{
			CartID: c.Value,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, items)
	}
}

// CartRemove takes out a product from the shopping cart.
func (s *Frontend) CartRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		q := chi.URLParam(r, "quantity")
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		quantity, err := strconv.Atoi(q)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		_, err = s.shoppingClient.Remove(ctx, &cart.RemoveRequest{
			CartID:    c.Value,
			ProductID: id,
			Quantity:  int64(quantity),
		})
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Successfully removed the product from the cart")
	}
}

// CartReset resets the cart to its default state.
func (s *Frontend) CartReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		_, err := s.shoppingClient.Reset(ctx, &cart.ResetRequest{
			CartID: c.Value,
		})
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, r, http.StatusOK, "Cart reseted")
	}
}

// CartSize returns the size of the shopping cart.
func (s *Frontend) CartSize() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("CID")
		ctx := r.Context()

		size, err := s.shoppingClient.Size(ctx, &cart.SizeRequest{
			CartID: c.Value,
		})
		if err != nil {
			response.Error(w, r, http.StatusNotFound, err)
			return
		}

		response.JSON(w, r, http.StatusOK, size)
	}
}
