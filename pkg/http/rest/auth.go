package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/internal/sanitize"
	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/auth"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:4000/login/oauth2/google",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET_ID"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	googleStateString = token.GenerateRunes(20)
)

// Login takes a user credentials and authenticates it.
func (s *API) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user auth.User
		ctx := r.Context()

		sID, err := r.Cookie("SID")
		if err == nil {
			sLog, _ := s.sessionClient.AlreadyLoggedIn(ctx, &auth.AlreadyLoggedInRequest{SessionID: sID.Value})
			if sLog.LoggedIn {
				sID.MaxAge = int(sLog.SessionLen)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		if err := validator.New().StructCtx(ctx, &user); err != nil {
			http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
			return
		}

		if err := sanitize.Normalize(&user.Email, &user.Password); err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		login, err := s.sessionClient.Login(ctx, &auth.LoginRequest{Email: user.Email, Password: user.Password})
		if err != nil {
			response.Error(w, http.StatusUnauthorized, err)
			return
		}

		adminEmails := viper.GetStringMap("admins.email")
		if _, ok := adminEmails[user.Email]; ok {
			admID := token.GenerateRunes(8)
			auth.SetCookie(w, "AID", admID, "/", int(login.SessionLen))
		}

		userID, err := token.GenerateFixedJWT(login.UserID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, errors.Wrap(err, "failed generating a jwt token"))
			return
		}

		// -SID- is the user session id
		auth.SetCookie(w, "SID", login.SessionID, "/", int(login.SessionLen))
		// -UID- used to deny users from making requests to other accounts
		auth.SetCookie(w, "UID", userID, "/", int(login.SessionLen))
		// -CID- used to identify which cart belongs to each user
		auth.SetCookie(w, "CID", login.CartID, "/", int(login.SessionLen))

		response.JSONText(w, http.StatusOK, "logged in")
	}
}

// Logout logs the user out from the session and removes cookies.
func (s *API) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sID, err := r.Cookie("SID")
		if err != nil {
			response.Error(w, http.StatusBadRequest, errors.New("error: you cannot log out without a session"))
			return
		}

		// Logout user from the session
		s.sessionClient.Logout(ctx, &auth.LogoutRequest{SessionID: sID.Value})

		if _, err := r.Cookie("AID"); err == nil {
			auth.DeleteCookie(w, "AID")
		}

		auth.DeleteCookie(w, "SID")
		auth.DeleteCookie(w, "UID")
		auth.DeleteCookie(w, "CID")

		response.JSONText(w, http.StatusOK, "logged out")
	}
}

// LoginGoogle redirects the user to the google oauth2.
func (s *API) LoginGoogle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := googleOauthConfig.AuthCodeURL(googleStateString)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// OAUTH2Google executes the oauth2 login with Google.
func (s *API) OAUTH2Google() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := userInfoGoogle(r.FormValue("state"), r.FormValue("code"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		res := fmt.Sprintf("Content: %s", content)

		response.JSONText(w, http.StatusOK, res)
	}
}

func userInfoGoogle(state, code string) ([]byte, error) {
	if state != googleStateString {
		return nil, errors.New("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, errors.Wrap(err, "code exchange failed")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting user info")
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading response body")
	}

	return contents, nil
}
