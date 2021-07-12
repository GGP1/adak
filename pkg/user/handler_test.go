package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/email"
	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/auth"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/user"
	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

var (
	handler     user.Handler
	userService user.Service
	cartService cart.Service
)

type msgResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type errResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func TestMain(m *testing.M) {
	poolMc, resourceMc, mc, err := test.RunMemcached()
	if err != nil {
		logger.Fatal(err)
	}
	poolPg, resourcePg, db, err := test.RunPostgres()
	if err != nil {
		logger.Fatal(err)
	}

	userService = user.NewService(db, mc)
	cartService = cart.NewService(db, mc)
	handler = user.NewHandler(true, userService, cartService, email.Emailer{}, mc)

	code := m.Run()

	if err := poolMc.Purge(resourceMc); err != nil {
		logger.Fatal(err)
	}
	if err := poolPg.Purge(resourcePg); err != nil {
		logger.Fatal(err)
	}

	os.Exit(code)
}

func TestCreateHandler(t *testing.T) {
	u := user.AddUser{
		Email:    "test@test.com",
		Username: "test",
		Password: "testing123",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(u)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", &buf)

	handler.Create()(rec, req)

	var response user.ListUser
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, u.Email, response.Email)
	assert.Equal(t, u.Username, response.Username)

	c, err := cartService.Get(context.Background(), response.CartID)
	assert.NoError(t, err)
	assert.Equal(t, response.CartID, c.ID)
	assert.Equal(t, int64(0), c.Total.Int64)
}

func TestDeleteHandler(t *testing.T) {
	u := user.AddUser{
		ID:       uuid.NewString(),
		CartID:   "test_delete_cart",
		Email:    "test_delete@test.com",
		Username: "test_delete",
		Password: "test_delete",
	}

	err := userService.Create(context.Background(), u)
	assert.NoError(t, err)

	rdb := test.StartRedis(t)

	session := auth.NewSession(nil, rdb, config.Session{}, true)
	mux := chi.NewRouter()
	mux.Delete("/{id}", handler.Delete(session))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/"+u.ID, nil)
	test.AddCookie(t, req, "UID", u.ID)
	test.AddCookie(t, req, "CID", u.CartID)

	mux.ServeHTTP(rec, req)

	var response msgResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, u.ID, response.Message)
}

func TestGetHandler(t *testing.T) {
	u := user.AddUser{
		ID:       uuid.NewString(),
		Email:    "test_get@test.com",
		Username: "test_get",
		Password: "test_get",
	}

	err := userService.Create(context.Background(), u)
	assert.NoError(t, err)

	mux := chi.NewRouter()
	mux.Get("/", handler.Get())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	mux.ServeHTTP(rec, req)

	var response []user.ListUser
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	if !strings.HasPrefix(response[0].Email, "test") { // Cannot determine which one of the users we are getting first
		t.Fatal("Invalid email")
	}
}

func TestGetByHandler(t *testing.T) {
	u := user.AddUser{
		ID:       uuid.NewString(),
		Email:    "test_getby@test.com",
		Username: "test_getby",
		Password: "test_getby",
	}

	err := userService.Create(context.Background(), u)
	assert.NoError(t, err)

	mux := chi.NewRouter()
	mux.Get("/email/{email}", handler.GetByEmail())
	mux.Get("/id/{id}", handler.GetByID())
	mux.Get("/username/{username}", handler.GetByUsername())

	t.Run("Email", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/email/"+u.Email, nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		var response user.ListUser
		err = json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, u.ID, response.ID)
	})

	t.Run("ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/id/"+u.ID, nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		var response user.ListUser
		err = json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, u.ID, response.ID)
	})

	t.Run("Username", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/username/"+u.Username, nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		var response user.ListUser
		err = json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, u.ID, response.ID)
	})
}

func TestSearchHandler(t *testing.T) {
	u := user.AddUser{
		ID:       uuid.NewString(),
		Email:    "test_search@test.com",
		Username: "test_search",
		Password: "test_search",
	}

	err := userService.Create(context.Background(), u)
	assert.NoError(t, err)

	mux := chi.NewRouter()
	mux.Get("/{query}", handler.Search())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/"+u.ID, nil) // Search by ID

	mux.ServeHTTP(rec, req)

	var response []user.ListUser
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, u.ID, response[0].ID)
}

func TestUpdateHandler(t *testing.T) {
	u := user.AddUser{
		ID:       uuid.NewString(),
		Email:    "test_update@test.com",
		Username: "test_update",
		Password: "test_update",
	}

	err := userService.Create(context.Background(), u)
	assert.NoError(t, err)

	mux := chi.NewRouter()
	mux.Put("/{id}", handler.Update())

	uptUser := user.UpdateUser{
		Username: "test_new_username",
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(uptUser)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/"+u.ID, &buf)
	test.AddCookie(t, req, "UID", u.ID)

	mux.ServeHTTP(rec, req)

	var response msgResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, u.ID, response.Message)
}
