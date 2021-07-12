package params

import (
	"context"
	"encoding/base64"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/validate"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

const maxResults = 50

// Object types
const (
	User obj = iota
	Shop
	Product
	Review
	Order
)

type obj uint8

// Cursor contains the values used for pagination.
type Cursor struct {
	// Used defines if the client used a cursor or not
	Used      bool
	CreatedAt time.Time
	ID        string
}

// Query contains the request parameters provided by the client.
//
// TODO: include order field to let clients change the objects' order (DESC/ASC)
type Query struct {
	Cursor Cursor
	Limit  string
}

// DecodeCursor decodes de cursor and returns both it and time
func DecodeCursor(encodedCursor string) (Cursor, error) {
	if encodedCursor == "" {
		return Cursor{Used: false}, nil
	}

	cursor, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return Cursor{}, errors.Wrap(err, "decoding cursor")
	}

	split := strings.Split(string(cursor), ",")
	if len(split) != 2 {
		return Cursor{}, errors.New("invalid cursor")
	}

	t, err := time.Parse(time.RFC3339Nano, split[0])
	if err != nil {
		return Cursor{}, errors.Wrap(err, "parsing time")
	}

	c := Cursor{
		Used:      true,
		CreatedAt: t,
		ID:        split[1],
	}
	return c, nil
}

// EncodeCursor encodes time and id with base64.
func EncodeCursor(t time.Time, id string) string {
	key := t.Format(time.RFC3339Nano) + "," + id
	return base64.StdEncoding.EncodeToString([]byte(key))
}

// ParseQuery returns the url params received after validating them.
func ParseQuery(rawQuery string, obj obj) (Query, error) {
	// Note: values.Get() retrieves only the first parameter, it's better to avoid accessing
	// the map manually, also validate the input to avoid HTTP parameter pollution.
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return Query{}, err
	}

	cursor, err := DecodeCursor(values.Get("cursor"))
	if err != nil {
		return Query{}, err
	}

	limit, err := parseInt(values.Get("limit"), "20", maxResults)
	if err != nil {
		return Query{}, errors.Wrap(err, "limit")
	}

	params := Query{
		Cursor: cursor,
		Limit:  limit,
	}
	return params, nil
}

// URLID returns the id parsed from the url.
func URLID(ctx context.Context) (string, error) {
	id := chi.URLParamFromCtx(ctx, "id")
	if err := validate.UUID(id); err != nil {
		return "", err
	}
	return id, nil
}

// split is like strings.Split but returns nil if the slice is empty
func split(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

// parseInt parses an integer from a url value and validates it.
//
// Return string as it will be used in a postgres query.
func parseInt(value, def string, max int) (string, error) {
	switch value {
	case "":
		return def, nil
	default:
		// Convert to integer to valiate it's a number as is lower than the maximum
		i, err := strconv.Atoi(value)
		if err != nil {
			return "", errors.Wrap(err, "invalid number")
		}
		if i < 0 {
			return def, nil
		}
		if max > 0 && i > max {
			return "", errors.Errorf("number provided (%d) exceeded maximum (%d)", i, max)
		}
		return value, nil
	}
}
