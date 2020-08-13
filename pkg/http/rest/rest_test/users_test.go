package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/storage"
	"github.com/badoux/checkmail"
)

func TestUsersHandler(t *testing.T) {
	t.Run("Add", add)
	t.Run("List", list)
}

func list(t *testing.T) {
	db, _, _ := storage.PostgresConnect()

	repo := *new(listing.Repository)
	users := handler.Users{DB: db}

	req := httptest.NewRequest("GET", "localhost:4000/users", nil)
	rec := httptest.NewRecorder()

	handler := users.Get(repo)
	handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	t.Log("Given the need to check the user handler.")
	{
		t.Logf("\t Test 0: When checking status code.")
		{
			if res.StatusCode != http.StatusOK {
				t.Errorf("\t%s\tShould return %v: got %v", failed, http.StatusOK, res.StatusCode)
			}
		}
	}
}

func add(t *testing.T) {
	users := []struct {
		firstname string
		lastname  string
		email     string
		password  string
	}{
		{firstname: "Test", lastname: "Ing", email: "testing@hotmail.com", password: "testing"},
		{firstname: "Error", lastname: "Test", email: "errortest@gmail.com", password: "errortest"},
	}

	t.Log("Given the need to test user adding.")
	{
		for i, user := range users {
			t.Logf("\tTest %d: When checking input validation.", i)
			{
				if user.firstname == "" {
					t.Errorf("\t%s\tShould enter Firstname", failed)
				}
				t.Logf("\t%s\tShould enter Firstname", succeed)

				if user.lastname == "" {
					t.Errorf("\t%s\tShould enter Lastname", failed)
				}
				t.Logf("\t%s\tShould enter Lastname", succeed)

				if user.email == "" {
					t.Errorf("\t%s\tShould enter Email", failed)
				}
				t.Logf("\t%s\tShould enter Email", succeed)

				if err := checkmail.ValidateFormat(user.email); err != nil {
					t.Errorf("\t%s\tShould be a valid Email", failed)
				}
				t.Logf("\t%s\tShould be a valid Email", succeed)

				if user.password == "" {
					t.Errorf("\t%s\tShould enter Password", failed)
				}
				t.Logf("\t%s\tShould enter Password", succeed)
			}
		}
	}
}
