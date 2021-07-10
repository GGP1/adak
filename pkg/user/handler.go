package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/cookie"
	"github.com/GGP1/adak/internal/email"
	"github.com/GGP1/adak/internal/params"
	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/internal/validate"
	"github.com/GGP1/adak/pkg/auth"
	"github.com/GGP1/adak/pkg/shopping/cart"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type cursorResponse struct {
	NextCursor string     `json:"next_cursor,omitempty"`
	Users      []ListUser `json:"users,omitempty"`
}

// Handler handles user endpoints.
type Handler struct {
	userService Service
	development bool
	emailer     email.Emailer
	cache       *memcache.Client
	cartService cart.Service
}

// NewHandler returns a new user handler.
func NewHandler(dev bool, userS Service, cartS cart.Service, emailer email.Emailer, cache *memcache.Client) Handler {
	return Handler{
		development: dev,
		userService: userS,
		cartService: cartS,
		emailer:     emailer,
		cache:       cache,
	}
}

// Create creates a new user and saves it.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user AddUser
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validate.Struct(ctx, user); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		if !h.development {
			confirmationCode := token.RandString(20)
			if err := h.emailer.SendValidation(ctx, user.Username, user.Email, confirmationCode); err != nil {
				response.Error(w, http.StatusInternalServerError, err)
				return
			}
		}

		// Set fields here to make testing easier and normalize inputs
		user.ID = token.RandString(32)
		user.CartID = token.RandString(31)
		user.Username = sanitize.Normalize(user.Username)
		user.Email = sanitize.Normalize(user.Email)
		user.CreatedAt = time.Now()

		if err := h.userService.Create(ctx, user); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		if err := h.cartService.Create(ctx, user.CartID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		user.Password = "" // Do not return password
		response.JSON(w, http.StatusCreated, user)
	}
}

// Delete removes a user.
func (h *Handler) Delete(s auth.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()
		cartID, err := cookie.GetValue(r, "CID")
		if err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		if err := token.CheckPermits(r, id); err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		if err := h.cartService.Delete(ctx, cartID); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if err := h.userService.Delete(ctx, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if err := s.Logout(ctx, w, r); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, id)
	}
}

// Get lists all the users.
func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		urlParams, err := params.ParseQuery(r.URL.RawQuery, params.Shop)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		users, err := h.userService.Get(ctx, urlParams)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		var nextCursor string
		if len(users) > 0 {
			nextCursor = params.EncodeCursor(users[len(users)-1].CreatedAt, users[len(users)-1].ID)
		}

		response.JSON(w, http.StatusOK, cursorResponse{
			NextCursor: nextCursor,
			Users:      users,
		})
	}
}

// GetByID lists the user with the id requested.
func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		item, err := h.cache.Get(id)
		if err == nil {
			response.EncodedJSON(w, item.Value)
			return
		}

		user, err := h.userService.GetByID(ctx, id)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSONAndCache(h.cache, w, id, user)
	}
}

// GetByEmail lists the user with the id requested.
func (h *Handler) GetByEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		ctx := r.Context()

		user, err := h.userService.GetByEmail(ctx, email)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, user)
	}
}

// GetByUsername lists the user with the id requested.
func (h *Handler) GetByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		ctx := r.Context()

		user, err := h.userService.GetByUsername(ctx, username)
		if err != nil {
			response.Error(w, http.StatusNotFound, err)
			return
		}

		response.JSON(w, http.StatusOK, user)
	}
}

// Search looks for the products with the given value.
func (h *Handler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := chi.URLParam(r, "query")
		ctx := r.Context()

		query = sanitize.Normalize(query)
		if strings.ContainsAny(query, ";-\\|@#~€¬<>_()[]}{¡^'") {
			response.Error(w, http.StatusBadRequest, errors.New("query contains invalid characters"))
			return
		}

		users, err := h.userService.Search(ctx, query)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, users)
	}
}

// Update updates the user with the given id.
func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UpdateUser
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		if err := token.CheckPermits(r, id); err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validate.Struct(ctx, user); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		if err := h.userService.Update(ctx, user, id); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSONText(w, http.StatusOK, id)
	}
}
