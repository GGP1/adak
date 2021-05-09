package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	cases := []struct {
		desc     string
		level    level
		message  string
		expected string
	}{
		{
			desc:     "INFO",
			level:    info,
			message:  "test info",
			expected: "testing/testing.go#1194 - INFO: test info\n",
		},
		{
			desc:     "DEBUG",
			level:    debug,
			message:  "test debug",
			expected: "testing/testing.go#1194 - DEBUG: test debug\n",
		},
		{
			desc:     "ERROR",
			level:    err,
			message:  "test error",
			expected: "testing/testing.go#1194 - ERROR: test error\n",
		},
		{
			desc:     "FATAL",
			level:    fatal,
			message:  "test fatal",
			expected: "testing/testing.go#1194 - FATAL: test fatal\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			var buf bytes.Buffer
			logger := New(true, false, &buf)

			logger.log(tc.level, tc.message)

			assert.Equal(t, tc.expected, buf.String())
		})
	}
}
