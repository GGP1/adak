package params

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecodeCursor(t *testing.T) {
	expected := Cursor{
		CreatedAt: time.Unix(51000, 0),
		ID:        "1568741",
		Used:      true,
	}
	str := expected.CreatedAt.Format(time.RFC3339Nano) + "," + expected.ID
	encodedCursor := base64.StdEncoding.EncodeToString([]byte(str))

	got, err := DecodeCursor(encodedCursor)
	assert.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestEncodeCursor(t *testing.T) {
	unix := time.Unix(51000, 0)
	id := "1568741"
	str := unix.Format(time.RFC3339Nano) + "," + id
	expected := base64.StdEncoding.EncodeToString([]byte(str))

	got := EncodeCursor(unix, id)
	assert.Equal(t, expected, got)
}

func TestParseQuery(t *testing.T) {
	createdAt := time.Unix(15000, 0) // nsec must be 0
	id := "1234567890"
	encodedCursor := EncodeCursor(createdAt, id)
	cases := []struct {
		desc     string
		obj      obj
		rawQuery string
		expected Query
	}{
		{
			desc:     "User",
			obj:      User,
			rawQuery: "cursor=" + encodedCursor + "&limit=20",
			expected: Query{
				Cursor: Cursor{
					Used:      true,
					CreatedAt: createdAt,
					ID:        id,
				},
				Limit: "20",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := ParseQuery(tc.rawQuery, tc.obj)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestParseInt(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		expected := "20"
		got, err := parseInt("20", "12", 50)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("Default", func(t *testing.T) {
		expected := "12"
		got, err := parseInt("", "12", 50)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("Invalid", func(t *testing.T) {
		_, err := parseInt("abc", "12", 50)
		assert.Error(t, err)
	})
	t.Run("Maximum exceeded", func(t *testing.T) {
		_, err := parseInt("20", "12", 15)
		assert.Error(t, err)
	})
}

func TestSplit(t *testing.T) {
	cases := []struct {
		desc     string
		expected []string
		input    string
	}{
		{
			desc:     "Non-nil",
			expected: []string{"name", "username", "email", "birth_date"},
			input:    "name,username,email,birth_date",
		},
		{
			desc:     "Nil",
			expected: nil,
			input:    "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := split(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}
