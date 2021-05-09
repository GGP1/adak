package auth

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockSession struct{}

func (s *mockSession) AlreadyLoggedIn(ctx context.Context, r *http.Request) bool {
	return false
}
func (s *mockSession) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, email, password string) error {
	return nil
}
func (s *mockSession) LoginOAuth(ctx context.Context, w http.ResponseWriter, r *http.Request, email string) error {
	return nil
}
func (s *mockSession) Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func TestLoginHandler(t *testing.T) {
	// Actually I should use the real session instead
	var session *mockSession
	rec := httptest.NewRecorder()

	body := bytes.NewBufferString(`
	{
		"email": "some-email@provider.com",
		"password": "123456789"
	}`)

	req, err := http.NewRequest("POST", "https://localhost:4000/login", body)
	if err != nil {
		t.Fatalf("Failed sending login request: %v", err)
	}

	hFunc := Login(session)
	hFunc.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected OK, got %s", res.Status)
	}
}

func TestLogoutHandler(t *testing.T) {
	// Actually we should use the real session instead
	var session *mockSession
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "https://localhost:4000/logout", nil)
	if err != nil {
		t.Fatalf("Failed sending logout request: %v", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "SID",
		Value: "",
	})

	hFunc := Logout(session)
	hFunc.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected OK, got %s", res.Status)
	}
}

func TestLoginGoogleHandler(t *testing.T) {
	var session *mockSession
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "https://localhost:4000/login/google", nil)
	if err != nil {
		t.Fatalf("Failed sending logout request: %v", err)
	}

	hFunc := LoginGoogle(session)
	hFunc.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Expected OK, got %s", res.Status)
	}
}
