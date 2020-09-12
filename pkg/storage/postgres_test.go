package storage_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	cfg "github.com/GGP1/palo/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Checkboxes
const (
	succeed = "\u2713"
	failed  = "\u2717"
)

func TestPostgres(t *testing.T) {
	c, err := cfg.New()
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name, c.Database.SSLMode)

	t.Log("Given the need to test database connection.")
	{
		t.Logf("\tTest 0:\tWhen checking the database connection.")
		{
			db, err := sqlx.Open("postgres", url)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to connect to the database: %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to connect to the database.", succeed)
			defer db.Close()

			t.Logf("\tTest 1:\tWhen checking the db variable is valid.")
			{
				if reflect.TypeOf(db) != reflect.TypeOf(&sqlx.DB{}) {
					t.Fatalf("\t%s\tShould be of type *sqlx.DB: %T", failed, db)
				}
				t.Logf("\t%s\tShould be of type *sqlx.DB.", succeed)
			}

			t.Logf("\tTest 2:\tWhen verifying that the connection is still alive.")
			{
				if err := db.Ping(); err != nil {
					t.Fatalf("\t%s\tShould be nil: %v", failed, err)
				}
				t.Logf("\t%s\tShould be nil.", succeed)
			}
		}
	}
}
