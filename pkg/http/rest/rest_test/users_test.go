package rest_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/user"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestUsersHandler(t *testing.T) {
	t.Run("Add", add)
	t.Run("List", list)
}

func add(t *testing.T) {
	users := []struct {
		firstname string
		lastname  string
		email     string
		password  string
	}{
		{firstname: "Test", lastname: "Ing", email: "testing@hotmail.com", password: "testing"},
		{firstname: "Add", lastname: "Test", email: "addTest@gmail.com", password: "addtest"},
	}

	for _, user := range users {
		if user.firstname == "" {
			t.Errorf("Should enter Firstname")
		}

		if user.lastname == "" {
			t.Errorf("Should enter Lastname")
		}

		if user.email == "" {
			t.Errorf("Should enter Email")
		}

		if err := validateEmail(user.email); err != nil {
			t.Errorf("Should be a valid Email")
		}

		if user.password == "" {
			t.Errorf("Should enter Password")
		}
	}
}

func list(t *testing.T) {
	c, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name, c.Database.SSLMode)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		t.Errorf("couldn't open the database: %v", err)
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

	if res.StatusCode != http.StatusOK {
		t.Errorf("Should return %v: got %v", http.StatusOK, res.StatusCode)
	}
}

func validateEmail(email string) error {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}
