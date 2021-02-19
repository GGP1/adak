package product

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func TestProductCreate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`INSERT INTO products`).
		WithArgs("ITlhtZhknGNFqFcDWZWYTNuCjXpwGA", 13, "San diego", "Dairy",
			"Dulce de leche", "Dulce de leche made in Argentina", 700, 6, 15, 100, time.Time{}, time.Time{})

	db := sqlx.NewDb(mockDB, "sqlmock")

	body := bytes.NewBufferString(`
	{
		"Shop_id": "ITlhtZhknGNFqFcDWZWYTNuCjXpwGA",
		"Stock": 13,
		"Brand": "San diego",
		"Category": "Dairy",
		"Type": "Dulce de leche",
		"Description": "Dulce de leche made in Argentina",
		"Weight": 700,
		"Discount": 6,
		"Taxes": 15,
		"Subtotal": 100
	}`)

	req, err := http.NewRequest("POST", "http://localhost:4000/products/create", body)
	if err != nil {
		t.Fatalf("Failed sending the request: %v", err)
	}

	rec := httptest.NewRecorder()

	service := NewService(*new(Repository), db)
	h := Handler{
		Service: service,
	}

	hf := h.Create()
	hf.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %s", res.Status)
	}
}
