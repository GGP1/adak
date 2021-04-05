package review

import (
	"context"
	"testing"
	"time"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/postgres"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var r = Review{
	ID:        "test",
	Stars:     5,
	Comment:   "Testing is awesome",
	UserID:    "test",
	ProductID: "test",
	ShopID:    "test",
}

var invalidR = Review{
	ID:        "invalid",
	Stars:     0,
	Comment:   "",
	UserID:    "non-existent",
	ProductID: "non-existent",
	ShopID:    "non-existent",
}

func NewTestService() (Service, func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conf, err := config.New()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed creating a new configuration")
	}

	db, err := postgres.Connect(ctx, &conf.Database)
	if err != nil {
		return nil, nil, err
	}

	service := NewService(db)
	return service, db.Close, nil
}

func TestService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()

	s, close, err := NewTestService()
	if err != nil {
		t.Error(err)
	}
	defer close()

	if err := s.Create(ctx, &r, r.UserID); err != nil {
		t.Fatalf("Failed creating a review: %v", err)
	}

	if err := s.Delete(ctx, r.ID); err != nil {
		t.Fatalf("Failed deleting the review: %v", err)
	}

	if err := s.Create(ctx, &r, r.UserID); err != nil {
		t.Fatalf("Failed creating a review: %v", err)
	}

	got, err := s.GetByID(ctx, r.ID)
	if err != nil {
		t.Errorf("Review not found: %v", err)
	}

	if got.ID != r.ID {
		t.Errorf("Expected a review with the id %q, got %q", r.ID, got.ID)
	}

	reviews, err := s.Get(ctx)
	if err != nil {
		t.Errorf("Get() failed: %v", err)
	}

	if len(reviews) == 0 {
		t.Errorf("Expected atleast one review, got %d", len(reviews))
	}
}

func TestInvalidService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()

	s, close, err := NewTestService()
	if err != nil {
		t.Error(err)
	}
	defer close()

	if err := s.Create(ctx, &invalidR, invalidR.UserID); err == nil {
		t.Fatalf("Expected Create() to fail but it didn't: %v", err)
	}

	if err := s.Delete(ctx, invalidR.ID); err == nil {
		t.Errorf("Expected Delete() to fail but it didn't: %v", err)
	}

	if _, err := s.GetByID(ctx, invalidR.ID); err == nil {
		t.Errorf("Expected GetByID() to fail but it didn't: %v", err)
	}
}
