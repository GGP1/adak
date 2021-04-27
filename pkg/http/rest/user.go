package rest

import (
	"encoding/json"
	"net/http"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/user"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// UserCreate creates a new user and saves it.
func (s *Frontend) UserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var addUser user.AddUser
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&addUser); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, &addUser); err != nil {
			http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
			return
		}

		if _, err := s.userClient.Create(ctx, &user.CreateRequest{User: &addUser}); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		response.HTMLText(w, http.StatusCreated, "Your account was successfully created.")
	}
}

// UserDelete removes a user.
func (s *Frontend) UserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uID, _ := r.Cookie("UID")
		ctx := r.Context()

		if err := token.CheckPermits(id, uID.Value); err != nil {
			response.Error(w, http.StatusUnauthorized, err)
			return
		}

		getByID, err := s.userClient.GetByID(ctx, &user.GetByIDRequest{ID: id})
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		_, err = s.shoppingClient.Delete(ctx, &cart.DeleteRequest{CartID: getByID.User.CartID})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		_, err = s.userClient.Delete(ctx, &user.DeleteRequest{ID: id})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, http.StatusOK, "User deleted successfully.")
		http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
	}
}

// UserGet lists all the users.
func (s *Frontend) UserGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		users, err := s.userClient.Get(ctx, &user.GetRequest{})
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, users)
	}
}

// UserGetByID lists the user with the id requested.
func (s *Frontend) UserGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		user, err := s.userClient.GetByID(ctx, &user.GetByIDRequest{ID: id})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, user)
	}
}

// UserSearch looks for the products with the given value.
func (s *Frontend) UserSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")
		ctx := r.Context()

		if err := sanitize.Normalize(&query); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		users, err := s.userClient.Search(ctx, &user.SearchRequest{Search: query})
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, users)
	}
}

// UserUpdate updates the user with the given id.
func (s *Frontend) UserUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateUser user.UpdateUser
		id := chi.URLParam(r, "id")
		uID, _ := r.Cookie("UID")
		ctx := r.Context()

		if err := token.CheckPermits(id, uID.Value); err != nil {
			response.Error(w, http.StatusUnauthorized, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, &updateUser); err != nil {
			http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
			return
		}

		_, err := s.userClient.Update(ctx, &user.UpdateRequest{User: &updateUser, ID: id})
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.HTMLText(w, http.StatusOK, "User updated successfully.")
	}
}
