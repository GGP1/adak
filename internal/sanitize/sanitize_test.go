package sanitize

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{in: "Dança", expected: "Danca"},
		{in: "Çomer", expected: "Comer"},
		{in: "úser", expected: "user"},
		{in: "ïd", expected: "id"},
		{in: "nÀmệ", expected: "nAme"},
	}

	for i, tc := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := Normalize(tc.in)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func BenchmarkNormalize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Normalize("BénçhmẬrkstrïng") // Maybe it requires a little more research to estimate a good average input
	}
}
