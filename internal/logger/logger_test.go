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
		prefix   string
		message  string
		expected string
	}{
		{
			desc:     "INFO",
			level:    Info,
			prefix:   "[TEST]",
			message:  "test info",
			expected: "[TEST] - INFO: test info\n",
		},
		{
			desc:     "DEBUG",
			level:    Debug,
			prefix:   "[ADAK]",
			message:  "test debug",
			expected: "[ADAK] - DEBUG: test debug\n",
		},
		{
			desc:     "ERROR",
			level:    Error,
			prefix:   "[ADAK]",
			message:  "test error",
			expected: "[ADAK] - ERROR: test error\n",
		},
		{
			desc:     "FATAL",
			level:    Fatal,
			prefix:   "[MONOLITHIC]",
			message:  "test fatal",
			expected: "[MONOLITHIC] - FATAL: test fatal\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var buf bytes.Buffer
			logger.Prefix = tc.prefix
			logger.ShowTimestamp = false
			logger.Out = &buf

			logger.log(tc.level, tc.message)

			assert.Equal(t, tc.expected, buf.String())
		})
	}
}
