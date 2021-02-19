package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GZIPCompress checks if the request accepts encoding and utilized gzip or proceed without compressing.
func GZIPCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gw := NewGZIPResponseWriter(w)
			gw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(gw, r)

			gw.Flush()
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GZIPReponseWriter is a response writer containing a GZIP writer in it
type GZIPReponseWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

// NewGZIPResponseWriter returns a new GZIPResponseWriter.
func NewGZIPResponseWriter(rw http.ResponseWriter) *GZIPReponseWriter {
	gw := gzip.NewWriter(rw)

	return &GZIPReponseWriter{w: rw, gw: gw}
}

// Header is implemented to satisfy the response writer interface.
func (g *GZIPReponseWriter) Header() http.Header {
	return g.w.Header()
}

// Write is implemented to satisfy the response writer interface.
func (g *GZIPReponseWriter) Write(d []byte) (int, error) {
	return g.gw.Write(d)
}

// WriteHeader is implemented to satisfy the response writer interface.
func (g *GZIPReponseWriter) WriteHeader(statuscode int) {
	g.w.WriteHeader(statuscode)
}

// Flush flushes any pending compressed data and closes the gzip writer.
func (g *GZIPReponseWriter) Flush() {
	g.gw.Flush()
	g.gw.Close()
}
