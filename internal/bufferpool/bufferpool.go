// Package bufferpool provides a pool of bytes buffers to avoid allocations and reusing them.
package bufferpool

import (
	"bytes"
	"sync"
)

var pool = &sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// Get returns a buffer from the pool.
func Get() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

// Put resets buf and puts it back to the pool.
func Put(buf *bytes.Buffer) {
	buf.Reset()
	pool.Put(buf)
}
