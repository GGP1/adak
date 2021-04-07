package tracking

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// Footprint takes user non-private information and generates a hash.
// The hash is unique for the visitor, not for the page.
func Footprint(r *http.Request, salt string) (string, error) {
	var sb strings.Builder

	ip, err := getUserIP(r)
	if err != nil {
		return "", err
	}

	sb.WriteString(r.Header.Get("User-Agent"))
	sb.WriteString(ip)
	sb.WriteString(fmt.Sprint(time.Now().UnixNano()))
	sb.WriteString(salt)
	hash := md5.New()

	if _, err := io.WriteString(hash, sb.String()); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func getUserIP(r *http.Request) (string, error) {
	forwarded := r.Header.Get("X-Forwarded-For")

	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0]), nil
	}

	forwarded = r.Header.Get("Forwarded")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		f := strings.Split(ips[0], ";")
		return strings.TrimSpace(f[0]), nil
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return ip, nil
}
