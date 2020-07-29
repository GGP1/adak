package tracking

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// Footprint returns a hash for given request and salt.
// The hash is unique for the visitor, not for the page.
func Footprint(r *http.Request, salt string) string {
	var sb strings.Builder

	sb.WriteString(r.Header.Get("User-Agent"))
	sb.WriteString(getIP(r))
	sb.WriteString(time.Now().UTC().Format("20060102"))
	sb.WriteString(salt)
	hash := md5.New()

	if _, err := io.WriteString(hash, sb.String()); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}

// Get user IP
func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")

	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// alternative header
	forwarded = r.Header.Get("Forwarded")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		f := strings.Split(ips[0], ";")
		return strings.TrimSpace(f[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}
