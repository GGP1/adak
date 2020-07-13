package handler

import (
	"errors"
	"net/http"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/email"
	"github.com/gorilla/mux"
)

// ValidateEmail is the email verification page
func ValidateEmail(pendingList, validatedList email.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var validated bool
		token := mux.Vars(r)["token"]

		pList, err := pendingList.Read()
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		for k, v := range pList {
			if v == token {
				err := validatedList.Add(k, v)
				if err != nil {
					response.Error(w, r, http.StatusInternalServerError, err)
				}
				validated = true
			}
		}

		if !validated {
			response.Error(w, r, http.StatusInternalServerError, errors.New("An error ocurred when validating your email"))
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You have successfully validated your email!")
	}
}
