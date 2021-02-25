package user

import (
	"context"
	"testing"
	"time"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/pkg/postgres"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var u = &AddUser{
	ID:       "test",
	CartID:   "test",
	Username: "test",
	Email:    "test@test.com",
	Password: "testing123",
}

var invalidU = &AddUser{
	ID:       "invalid",
	CartID:   "non-existent",
	Username: "non-existent",
	Email:    "invalid",
	Password: "inv",
}

func NewUserService(t *testing.T) (context.Context, func() error, Service) {
	t.Helper()
	logger.Log.Disable()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)

	conf, err := config.New()
	assert.NoError(t, err)

	conf.Database = config.Database{
		Host:     "localhost",
		Port:     "6000",
		Username: "test",
		Password: "test",
		Name:     "test",
		SSLMode:  "disable",
	}

	db, err := postgres.Connect(ctx, &conf.Database)
	assert.NoError(t, err)

	repo := *new(Repository)
	service := NewService(repo, db)

	closure := func() error {
		cancel()
		return db.Close()
	}

	return ctx, closure, service
}

func TestCreate(t *testing.T) {
	ctx, close, s := NewUserService(t)
	defer close()

	assert.NoError(t, s.Create(ctx, u))
}

func TestDelete(t *testing.T) {
	ctx, close, s := NewUserService(t)
	defer close()

	assert.NoError(t, s.Delete(ctx, u.ID))
}

func TestGet(t *testing.T) {
	ctx, close, s := NewUserService(t)
	defer close()

	_, err := s.Get(ctx)
	assert.NoError(t, err)
}

func TestGetByID(t *testing.T) {
	ctx, close, s := NewUserService(t)
	defer close()

	user, err := s.GetByID(ctx, u.ID)
	assert.NoError(t, err)

	assert.Equal(t, u.ID, user.ID)
}

func TestGetByEmail(t *testing.T) {
	ctx, close, s := NewUserService(t)
	defer close()

	user, err := s.GetByEmail(ctx, u.Email)
	if err != nil {
		t.Errorf("Failed creating a user: %v", err)
	}

	assert.Equal(t, u.Email, user.Email)
}

func TestGetByUsername(t *testing.T) {
	ctx, close, s := NewUserService(t)
	defer close()

	user, err := s.GetByEmail(ctx, u.Email)
	assert.NoError(t, err)

	assert.Equal(t, u.Username, user.Username)
}

func TestUpdate(t *testing.T) {
	ctx, close, s := NewUserService(t)
	defer close()

	username := "newUsername"
	assert.NoError(t, s.Update(ctx, &UpdateUser{Username: username}, u.ID))

	uptUser, err := s.GetByUsername(ctx, username)
	assert.NoError(t, err)

	assert.Equal(t, username, uptUser.Username)
}

func TestSearch(t *testing.T) {
	ctx, close, s := NewUserService(t)
	defer close()

	users, err := s.Search(ctx, u.ID)
	assert.NoError(t, err)

	var found bool
	for _, us := range users {
		if us.ID == u.ID {
			found = true
			break
		}
	}

	assert.Equal(t, found, true)
}

func TestInvalidService(t *testing.T) {
	ctx, close, s := NewUserService(t)
	defer close()

	assert.Error(t, s.Create(ctx, invalidU))
	assert.Error(t, s.Delete(ctx, invalidU.ID))

	_, err := s.GetByID(ctx, invalidU.ID)
	assert.Error(t, err)

	username := "invalidUsername"
	assert.Error(t, s.Update(ctx, &UpdateUser{Username: username}, u.ID))

	uptUser, err := s.GetByUsername(ctx, username)
	assert.Error(t, err)

	assert.NotEqual(t, uptUser.Username, username)

	_, err = s.Search(ctx, invalidU.ID)
	assert.Error(t, err)
}
