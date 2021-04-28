package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	logger := New()

	testCases := []struct {
		desc     string
		level    Level
		message  string
		expected string
	}{
		{
			desc:     "INFO",
			level:    Info,
			message:  "test info",
			expected: "testing/testing.go#1194 - INFO: test info\n",
		},
		{
			desc:     "DEBUG",
			level:    Debug,
			message:  "test debug",
			expected: "testing/testing.go#1194 - DEBUG: test debug\n",
		},
		{
			desc:     "ERROR",
			level:    Error,
			message:  "test error",
			expected: "testing/testing.go#1194 - ERROR: test error\n",
		},
		{
			desc:     "FATAL",
			level:    Fatal,
			message:  "test fatal",
			expected: "testing/testing.go#1194 - FATAL: test fatal\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var buf bytes.Buffer
			logger.showTimestamp = false
			logger.out = &buf

			logger.log(tc.level, tc.message)

			assert.Equal(t, tc.expected, buf.String())
		})
	}
}
