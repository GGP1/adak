package validate

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type value struct {
	ID          string `validate:"required"`
	Events      []string
	Email       string `validate:"required,email"`
	Username    string `validate:"required,min=4"`
	Age         uint16 `validate:"required,gte=0,lte=128"`
	Description string `validate:"max=144"`
	Host        bool
	Gamertag    string
}

func TestStruct(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		v := value{
			ID:       "F53tPeXayYG6Bs",
			Email:    "adak@test.com",
			Username: "adak",
			Age:      27,
		}
		assert.NoError(t, Struct(context.Background(), v))
	})

	t.Run("Validation error", func(t *testing.T) {
		v := value{
			ID:       "VTizDf64ylNQmF",
			Email:    "adaktest.com",
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
