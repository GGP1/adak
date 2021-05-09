package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/auth"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
)

var (
	session auth.Session
	db      *sqlx.DB
	rdb     *redis.Client
	email   = "test_auth_session@test.com"
)

func TestMain(m *testing.M) {
	config := config.Session{
		Attempts: 1,
		Delay:    0,
	}
	pgPool, pgResource, sqlxDB, err := test.RunPostgres()
	if err != nil {
		logger.Fatal(err)
	}

	rdbPool, rdbResource, redisDB, err := test.RunRedis()
	if err != nil {
		logger.Fatal(err)
	}
	db = sqlxDB
	rdb = redisDB

	session = auth.NewSession(db, rdb, config, true)
	if err := createUser(context.Background()); err != nil {
		logger.Fatal(err)
	}

	code := m.Run()

	if err := pgPool.Purge(pgResource); err != nil {
		logger.Fatal(err)
	}
	if err := rdbPool.Purge(rdbResource); err != nil {
		logger.Fatal(err)
	}

	os.Exit(code)
}

func TestAlreadyLoggedIn(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		ctx := context.Background()
		salt := "0123456789111213"
		sID := "1:" + salt
		err := rdb.Set(ctx, sID, salt, 0).Err()
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		test.AddCookie(t, r, "SID", sID)
		got := session.AlreadyLoggedIn(ctx, r)
		assert.Equal(t, true, got)
	})

	t.Run("False", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		got := session.AlreadyLoggedIn(context.Background(), r)
		assert.Equal(t, false, got)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Standard", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		err := session.Login(context.Background(), rec, req, email, "password")
		assert.NoError(t, err)

		cookies := rec.Result().Cookies()
		assert.Equal(t, 3, len(cookies))
		assert.Equal(t, "SID", cookies[0].Name)
		assert.Equal(t, "UID", cookies[1].Name)
		assert.Equal(t, "CID", cookies[2].Name)
	})

	t.Run("OAuth", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		err := session.LoginOAuth(context.Background(), rec, req, email)
		assert.NoError(t, err)

		cookies := rec.Result().Cookies()
		assert.Equal(t, 3, len(cookies))
		assert.Equal(t, "SID", cookies[0].Name)
		assert.Equal(t, "UID", cookies[1].Name)
		assert.Equal(t, "CID", cookies[2].Name)
	})
}

func TestLogout(t *testing.T) {
	ctx := context.Background()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	err := session.Logout(ctx, rec, req)
	assert.NoError(t, err)

	cookies := rec.Result().Cookies()
	// They are not deleted by the recorder
	assert.Equal(t, 3, len(cookies))
	assert.Equal(t, "", cookies[0].Value)
	assert.Equal(t, "", cookies[1].Value)
	assert.Equal(t, "", cookies[2].Value)
}

func createUser(ctx context.Context) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	q := `INSERT INTO users
	(id, cart_id, username, email, password, verified_email)
	VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.ExecContext(ctx, q, "1", "2", "username", email, hash, true)
	if err != nil {
		return err
	}

	return nil
}
