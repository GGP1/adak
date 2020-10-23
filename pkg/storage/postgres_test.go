package storage_test

import (
	"fmt"
	"testing"

	"github.com/GGP1/palo/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Checkbox
const failed = "\u2717"

func TestPostgres(t *testing.T) {
	c, err := config.New()
	if err != nil {
		t.Fatalf("Creating config failed: %v", err)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name, c.Database.SSLMode)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		t.Fatalf("\t%s\tFailed connecting to the database: %v", failed, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("\t%s\tPing to the database should be nil: %v", failed, err)
	}
}
