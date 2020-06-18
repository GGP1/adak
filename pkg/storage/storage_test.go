package storage_test

import (
	"testing"

	db "github.com/GGP1/palo/internal/utils/database"
	_ "github.com/lib/pq"
)

// Checkboxes
const (
	succeed = "\u2713"
	failed  = "\u2717"
)

func TestDatabase(t *testing.T) {
	t.Run("database", database)
}

func database(t *testing.T) {
	t.Log("Given the need to test database connection.")
	{
		t.Logf("\tTest 0:\tWhen checking the database connection.")
		{
			db, err := db.Connect(db.URL)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to connect to the database: %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to connect to the database.", succeed)

			defer db.Close()
		}
	}
}
