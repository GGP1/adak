package bufferpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferPool(t *testing.T) {
	text := "bufferpool test"
	buf := Get()
	buf.WriteString(text)
	assert.Equal(t, buf.String(), text)
	Put(buf)
	assert.Equal(t, buf.Len(), 0)
}
