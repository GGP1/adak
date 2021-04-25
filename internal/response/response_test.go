package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/GGP1/adak/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestEncodedJSON(t *testing.T) {
	expected := []byte("test")
	rec := httptest.NewRecorder()
	EncodedJSON(rec, []byte("test"))
	res := rec.Result()

	var buf bytes.Buffer
	_, err := buf.ReadFrom(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, expected, buf.Bytes())
}

func TestError(t *testing.T) {
	expectedHeaderCT := "application/json; charset=UTF-8"
	expectedStatus := 404
	expectedText := "{\"status\":404,\"error\":\"test\"}\n"

	rec := httptest.NewRecorder()
	Error(rec, http.StatusNotFound, errors.New("test"))
	res := rec.Result()

	gotHeaderCT := res.Header.Get("Content-Type")
	assert.Equal(t, expectedHeaderCT, gotHeaderCT)

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

func TestJSONAndCache(t *testing.T) {
	mc := test.StartMemcached(t)
	expectedHeader := "application/json; charset=UTF-8"
	expectedStatus := 200
	expectedRes := "\"test\"\n"
	key := "test_cache"
	value := "test"

	rec := httptest.NewRecorder()
	JSONAndCache(mc, rec, key, value)
	res := rec.Result()

	gotHeader := res.Header.Get("Content-Type")
	assert.Equal(t, expectedHeader, gotHeader)

	gotStatus := res.StatusCode
	assert.Equal(t, expectedStatus, gotStatus)

	var resContent bytes.Buffer
	_, err := resContent.ReadFrom(res.Body)
	assert.NoError(t, err, "Failed reading response body")
	assert.Equal(t, expectedRes, resContent.String())

	item, err := mc.Get(key)
	assert.NoError(t, err)

	var cacheContent bytes.Buffer
	err = json.NewEncoder(&cacheContent).Encode(value)
	assert.NoError(t, err)
	assert.Equal(t, cacheContent.Bytes(), item.Value)
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

func TestPNGAndCache(t *testing.T) {
	mc := test.StartMemcached(t)
	testImage := image.NewRGBA(image.Rect(15, 15, 30, 30))
	key := "test-png-and-cache"

	var imageBuf bytes.Buffer
	if err := png.Encode(&imageBuf, testImage); err != nil {
		t.Fatalf("Failed encoding image: %v", err)
	}

	expectedHeaderCT := "image/png"
	expectedStatus := 200
	expectedImage := imageBuf.Bytes()
	expectedHeaderCL := strconv.Itoa(imageBuf.Len())

	rec := httptest.NewRecorder()
	PNGAndCache(rec, mc, key, testImage)
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

	item, err := mc.Get(key)
	assert.NoError(t, err)

	var cacheContent bytes.Buffer
	err = png.Encode(&cacheContent, testImage)
	assert.NoError(t, err)
	assert.Equal(t, cacheContent.Bytes(), item.Value)
}
