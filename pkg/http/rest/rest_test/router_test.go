package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/email"
	"github.com/GGP1/palo/pkg/http/rest"
	"github.com/GGP1/palo/pkg/http/rest/handler"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/storage"
	"github.com/GGP1/palo/pkg/updating"
)

// Checkboxes
const (
	succeed = "\u2713"
	failed  = "\u2717"
)

// Fix
func TestRouting(t *testing.T) {
	db, close, err := storage.NewDatabase()
	if err != nil {
		t.Fatal("Database failed connecting")
	}
	defer close()
	// Repos
	addingRepo := *new(adding.Repository)
	deletingRepo := *new(deleting.Repository)
	listingRepo := *new(listing.Repository)
	updatingRepo := *new(updating.Repository)
	sessionRepo := *new(handler.AuthRepository)
	emailRepo := *new(email.Repository)

	// Services
	adder := adding.NewService(addingRepo)
	deleter := deleting.NewService(deletingRepo)
	lister := listing.NewService(listingRepo)
	updater := updating.NewService(updatingRepo)
	// -- Session--
	session := handler.NewSession(sessionRepo)
	// -- Email lists --
	pendingList := email.NewPendingList(db, emailRepo)
	validatedList := email.NewValidatedList(db, emailRepo)

	srv := httptest.NewServer(rest.NewRouter(db, adder, deleter, lister, updater, session, pendingList, validatedList))
	defer srv.Close()

	t.Log("Given the need to test the router.")
	{
		t.Logf("\tTest 0: When checking GET request.")
		{
			res, err := http.Get("http://localhost:4000/")
			if err != nil {
				t.Errorf("\t%s\tShould return a response: %v", failed, err)
			}
			t.Logf("\t%s\tShould return a response.", succeed)

			if res.StatusCode != http.StatusOK {
				t.Errorf("\t%s\tShould be status OK: got %v", failed, res.StatusCode)
			}
			t.Logf("\t%s\tShould be status OK.", succeed)
		}
	}

}
