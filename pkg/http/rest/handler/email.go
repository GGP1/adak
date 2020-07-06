package handler

import (
	"errors"
	"net/http"

	"github.com/GGP1/palo/internal/email"
	"github.com/GGP1/palo/internal/response"
	"github.com/gorilla/mux"
)

// ValidateEmail is the email verification page
func ValidateEmail(pendigList *email.PendingList, validatedList *email.ValidatedList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var validated bool = false

		token := mux.Vars(r)["token"]

		for k, v := range pendigList.UserList {
			if v == token {
				// k = pendigList[user.Email]
				validatedList.Add(k, token)

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
