package handler

import (
	"errors"
	"net/http"

	"github.com/GGP1/palo/internal/email"
	"github.com/GGP1/palo/internal/response"
	"github.com/gorilla/mux"
)

// Email verification page
func Verify(pendigList *email.PendingList, verifiedList *email.VerifiedList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var validated bool = false

		token := mux.Vars(r)["token"]

		for k, v := range pendigList.UserList {
			if v == token {
				// k = pendigList[user.Email]
				verifiedList.Add(k, token)

				validated = true
			}
		}

		if !validated {
			response.Error(w, r, http.StatusInternalServerError, errors.New("An error ocurred when validating your email"))
			return
		}

		response.HTMLText(w, r, http.StatusOK, "You have successfully confirmed your email!")
	}
}
