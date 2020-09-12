package rest_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/GGP1/palo/internal/config"
	"github.com/GGP1/palo/pkg/user"
	"github.com/jmoiron/sqlx"
)

func TestUsersHandler(t *testing.T) {
	t.Run("Add", add)
	t.Run("List", list)
}

func list(t *testing.T) {
	c, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name, c.Database.SSLMode)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		t.Error("couldn't open the database")
	}

	repo := *new(user.Repository)
	service := user.NewService(repo, db)

	req := httptest.NewRequest("GET", "localhost:4000/users", nil)
	rec := httptest.NewRecorder()

	handler := user.Handler{Service: service}

	handle := handler.Get()
	handle(rec, req)

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

				if err := validateEmail(user.email); err != nil {
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

func validateEmail(email string) error {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}
