package rest_test

import (
	"testing"
)

func TestHome(t *testing.T) {
	// req := httptest.NewRequest("GET", "localhost:4000/", nil)
	// rec := httptest.NewRecorder()

	// handler.Home().ServeHTTP(rec, req)

	// res := rec.Result()
	// defer res.Body.Close()
	// t.Log("Given the need to check home handler.")
	// {
	// 	t.Logf("\tTest 0: When checking the response status.")
	// 	{
	// 		if res.StatusCode != http.StatusOK {
	// 			t.Errorf("\t%s\tShould be status OK: got %v", failed, res.StatusCode)
	// 		}
	// 		t.Logf("\t%s\tShould be status OK.", succeed)
	// 	}
	// 	t.Logf("\tTest 1: When checking the body.")
	// 	{
	// 		b, err := ioutil.ReadAll(res.Body)
	// 		if err != nil {
	// 			t.Errorf("\t%s\tShould read response body: %v", failed, err)
	// 		}
	// 		t.Logf("\t%s\tShould read response body: %v", succeed, string(b))
	// 	}
	// }
}
