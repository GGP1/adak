package validate

import (
	"context"
	"fmt"
	"testing"

	"github.com/GGP1/adak/internal/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type value struct {
	ID          string `validate:"uuid4_rfc4122"`
	Events      []string
	Email       string `validate:"required,email"`
	Username    string `validate:"required,min=4"`
	Age         uint16 `validate:"required,gte=0,lte=128"`
	Description string `validate:"max=144"`
	Host        bool
	Gamertag    string
}

func TestSearchQuery(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		query := "searching for values"
		err := SearchQuery(query)
		assert.NoError(t, err)
	})
	t.Run("Invalid", func(t *testing.T) {
		query := "'; DROP TABLE users; --"
		err := SearchQuery(query)
		assert.Error(t, err)
	})
}

func TestStruct(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		v := value{
			ID:       uuid.NewString(),
			Email:    "adak@test.com",
			Username: "adak",
			Age:      27,
		}
		assert.NoError(t, Struct(context.Background(), v))
	})

	t.Run("Validation error", func(t *testing.T) {
		v := value{
			ID:       "VTizDf64ylNQmF",
			Email:    "adak@test.com",
			Username: "mav",
			Age:      130,
		}
		assert.Error(t, Struct(context.Background(), v), "Expected an error and got nil")
	})

	t.Run("Invalid req body", func(t *testing.T) {
		err := Struct(context.Background(), nil)
		assert.Error(t, err, "Expected an error and got nil")
		assert.Equal(t, err, fmt.Errorf("invalid request body: validator: (nil)"))
	})
}

func TestUUID(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		id := uuid.NewString()
		err := UUID(id)
		assert.NoError(t, err)
	})

	t.Run("Invalid", func(t *testing.T) {
		id := token.RandString(36)
		err := UUID(id)
		assert.Error(t, err)
	})
}
