package response

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	expectedHeaderCT := "application/json; charset=UTF-8"
	expectedHeaderXCTO := "nosniff"
	expectedStatus := 404
	expectedText := "{\"status\":404,\"error\":\"test\"}\n"

	rec := httptest.NewRecorder()
	Error(rec, http.StatusNotFound, errors.New("test"))
	res := rec.Result()

	gotHeaderCT := res.Header.Get("Content-Type")
	assert.Equal(t, expectedHeaderCT, gotHeaderCT)

	gotHeaderXCTO := res.Header.Get("X-Content-Type-Options")
	assert.Equal(t, expectedHeaderXCTO, gotHeaderXCTO)

	gotStatus := res.StatusCode
	assert.Equal(t, expectedStatus, gotStatus)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(res.Body)
	assert.NoError(t, err, "Failed reading response body")

	assert.Equal(t, expectedText, buf.String())
}

func TestHTMLText(t *testing.T) {
	expectedHeader := "text/html; charset=UTF-8"
	expectedStatus := 200
	expectedText := "test\n"

	rec := httptest.NewRecorder()
	HTMLText(rec, http.StatusOK, "test")
	res := rec.Result()

	gotHeader := res.Header.Get("Content-Type")
	assert.Equal(t, expectedHeader, gotHeader)

	gotStatus := res.StatusCode
	assert.Equal(t, expectedStatus, gotStatus)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(res.Body)
	assert.NoError(t, err, "Failed reading response body")

	assert.Equal(t, expectedText, buf.String())
}

func TestJSON(t *testing.T) {
	expectedHeader := "application/json; charset=UTF-8"
	expectedStatus := 201
	expectedText := "\"test\"\n"

	rec := httptest.NewRecorder()
	JSON(rec, http.StatusCreated, "test")
	res := rec.Result()

	gotHeader := res.Header.Get("Content-Type")
	assert.Equal(t, expectedHeader, gotHeader)

	gotStatus := res.StatusCode
	assert.Equal(t, expectedStatus, gotStatus)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(res.Body)
	assert.NoError(t, err, "Failed reading response body")

	assert.Equal(t, expectedText, buf.String())
}

func TestJSONText(t *testing.T) {
	expectedHeader := "application/json; charset=UTF-8"
	expectedStatus := 200
	expectedText := "{\"message\":\"test\",\"code\":200}\n"

	rec := httptest.NewRecorder()
	JSONText(rec, http.StatusOK, "test")
	res := rec.Result()

	gotHeader := res.Header.Get("Content-Type")
	assert.Equal(t, expectedHeader, gotHeader)

	gotStatus := res.StatusCode
	assert.Equal(t, expectedStatus, gotStatus)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(res.Body)
	assert.NoError(t, err, "Failed reading response body")

	assert.Equal(t, expectedText, buf.String())
}

func TestPNG(t *testing.T) {
	var imageBuf bytes.Buffer
	testImage := image.NewRGBA(image.Rect(15, 15, 30, 30))
	if err := png.Encode(&imageBuf, testImage); err != nil {
		t.Fatalf("Failed encoding image: %v", err)
	}

	expectedHeaderCT := "image/png"
	expectedStatus := 200
	expectedImage := imageBuf.Bytes()
	expectedHeaderCL := strconv.Itoa(imageBuf.Len())

	rec := httptest.NewRecorder()
	PNG(rec, http.StatusOK, testImage)
	res := rec.Result()

	gotHeaderCT := res.Header.Get("Content-Type")
	assert.Equal(t, expectedHeaderCT, gotHeaderCT)

	gotHeaderCL := res.Header.Get("Content-Length")
	assert.Equal(t, expectedHeaderCL, gotHeaderCL)

	gotStatus := res.StatusCode
	assert.Equal(t, expectedStatus, gotStatus)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(res.Body)
	assert.NoError(t, err, "Failed reading response body")

	assert.Equal(t, expectedImage, buf.Bytes())
}
