package auth

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/token"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleState  = token.RandString(20)
	googleConfig = &oauth2.Config{
		RedirectURL: "http://localhost:4000/login/oauth2/google",
		Scopes:      []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:    google.Endpoint,
	}
)

// There is no need to validate this data as it comes directly from Google.
type oauthRes struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

// Login takes a user credentials and authenticates it.
func Login(s Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.AlreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var auth UserAuth
		ctx := r.Context()

		if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, auth); err != nil {
			response.Error(w, http.StatusBadRequest, err.(validator.ValidationErrors))
			return
		}

		auth.Email = sanitize.Normalize(auth.Email)
		auth.Password = sanitize.Normalize(auth.Password)

		if err := s.Login(ctx, w, r, auth.Email, auth.Password); err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		response.JSONText(w, http.StatusOK, "logged in")
	}
}

// Logout logs the user out from the session and removes cookies.
func Logout(s Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logout user from the session and delete cookies
		s.Logout(w, r)

		response.JSONText(w, http.StatusOK, "logged out")
	}
}

// LoginGoogle redirects the user to the google oauth2.
func LoginGoogle(s Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.AlreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		googleConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
		googleConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
		url := googleConfig.AuthCodeURL(googleState)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// OAuth2Google executes the oauth2 login with Google.
func OAuth2Google(s Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var oauth oauthRes
		ctx := r.Context()

		res, err := userInfoGoogle(r.FormValue("state"), r.FormValue("code"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer res.Body.Close()

		if err := json.NewDecoder(res.Body).Decode(&oauth); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		if err := s.LoginOAuth(ctx, w, r, oauth.Email); err != nil {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		response.JSONText(w, http.StatusOK, "logged in")
	}
}

func userInfoGoogle(state, code string) (*http.Response, error) {
	if state != googleState {
		return nil, errors.New("invalid OAuth state")
	}

	token, err := googleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, errors.Wrap(err, "code exchange failed")
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting user info")
	}

	return res, nil
}
