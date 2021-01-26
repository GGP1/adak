package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		mustBeEqual bool
	}{
		{
			desc:        "Do not normalize",
			input:       "test",
			mustBeEqual: true,
		},
		{
			desc:        "Normalize",
			input:       "tëst",
			mustBeEqual: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Create a copy to compare
			copy := tc.input
			err := Normalize(&tc.input)
			assert.NoError(t, err)

			if tc.mustBeEqual {
				assert.Equal(t, copy, tc.input)
			} else {
				assert.NotEqual(t, copy, tc.input)
			}
		})
	}
}
