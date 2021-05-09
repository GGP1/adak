package user_test

import (
	"context"
	"testing"

	"github.com/GGP1/adak/internal/logger"
	"github.com/GGP1/adak/internal/test"
	"github.com/GGP1/adak/pkg/user"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var u = user.AddUser{
	ID:       "test",
	CartID:   "test",
	Username: "test",
	Email:    "test@test.com",
	Password: "testing123",
	IsAdmin:  true,
}

// TestMain failed when creating the user service.
func NewUserService(t *testing.T) (context.Context, user.Service) {
	t.Helper()
	logger.Disable()
	ctx, cancel := context.WithCancel(context.Background())

	db := test.StartPostgres(t)
	mc := test.StartMemcached(t)
	service := user.NewService(db, mc)

	viper.Set("admins", []string{"test@test.com"})

	t.Cleanup(func() {
		cancel()
	})

	return ctx, service
}

func TestUserService(t *testing.T) {
	ctx, s := NewUserService(t)

	t.Run("Create", create(ctx, s))
	t.Run("Get", get(ctx, s))
	t.Run("Get by id", getByID(ctx, s))
	t.Run("Get by email", getByEmail(ctx, s))
	t.Run("Get by username", getByUsername(ctx, s))
	t.Run("Is admin", isAdmin(ctx, s))
	t.Run("Update", update(ctx, s))
	t.Run("Search", search(ctx, s))
	t.Run("Delete", delete(ctx, s))
}

func create(ctx context.Context, s user.Service) func(t *testing.T) {
	return func(t *testing.T) {
		assert.NoError(t, s.Create(ctx, u))

		user, err := s.GetByID(ctx, u.ID)
		assert.NoError(t, err)

		assert.Equal(t, u.Username, user.Username)
	}
}

func delete(ctx context.Context, s user.Service) func(t *testing.T) {
	return func(t *testing.T) {
		assert.NoError(t, s.Delete(ctx, u.ID))

		user, err := s.GetByID(ctx, u.ID)
		assert.NoError(t, err)

		assert.Equal(t, "", user.ID)
	}
}

func get(ctx context.Context, s user.Service) func(t *testing.T) {
	return func(t *testing.T) {
		users, err := s.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, u.Email, users[0].Email)
	}
}

func getByID(ctx context.Context, s user.Service) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := s.GetByID(ctx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, u.ID, user.ID)
	}
}

func getByEmail(ctx context.Context, s user.Service) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := s.GetByEmail(ctx, u.Email)
		assert.NoError(t, err)
		assert.Equal(t, u.Email, user.Email)
	}
}

func getByUsername(ctx context.Context, s user.Service) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := s.GetByUsername(ctx, u.Username)
		assert.NoError(t, err)

		assert.Equal(t, u.Username, user.Username)
	}
}

func isAdmin(ctx context.Context, s user.Service) func(t *testing.T) {
	return func(t *testing.T) {
		isAdmin, err := s.IsAdmin(ctx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, u.IsAdmin, isAdmin)
	}
}

func update(ctx context.Context, s user.Service) func(t *testing.T) {
	return func(t *testing.T) {
		username := "newUsername"
		assert.NoError(t, s.Update(ctx, user.UpdateUser{Username: username}, u.ID))

		uptUser, err := s.GetByUsername(ctx, username)
		assert.NoError(t, err)

		assert.Equal(t, username, uptUser.Username)
	}
}

func search(ctx context.Context, s user.Service) func(t *testing.T) {
	return func(t *testing.T) {
		users, err := s.Search(ctx, u.ID)
		assert.NoError(t, err)

		var found bool
		for _, us := range users {
			if us.ID == u.ID {
				found = true
				break
			}
		}
		assert.Equal(t, true, found)
	}
}
