package tracking

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/bufferpool"
)

// Headers and corresponding parser to look up the real client IP.
// They will be check in order, the first non-empty one will be picked,
// or else the remote address is selected.
// CF-Connecting-IP is a header added by Cloudflare: https://support.cloudflare.com/hc/en-us/articles/206776727-What-is-True-Client-IP-
var ipHeaders = []ipHeader{
	{"CF-Connecting-IP", parseXForwardedForHeader},
	{"X-Forwarded-For", parseXForwardedForHeader},
	{"Forwarded", parseForwardedHeader},
	{"X-Real-IP", parseXRealIPHeader},
}

type ipHeader struct {
	header string
	parser func(string) string
}

// Footprint takes user non-private information and generates a hash.
// The hash is unique for the visitor, not for the page.
func Footprint(r *http.Request, salt string) (string, error) {
	ip := GetUserIP(r)
	buf := bufferpool.Get()
	buf.WriteString(r.Header.Get("User-Agent"))
	buf.WriteString(ip)
	buf.WriteString(fmt.Sprintf("%d", time.Now().UnixNano()))
	buf.WriteString(salt)

	hash := md5.New()
	// md5 write method returns always nil
	hash.Write(buf.Bytes())
	bufferpool.Put(buf)

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetUserIP returns the IP from given request.
// It will try to extract the real client IP from headers if possible.
func GetUserIP(r *http.Request) string {
	ip := r.RemoteAddr

	for _, header := range ipHeaders {
		value := r.Header.Get(header.header)

		if value != "" {
			parsedIP := header.parser(value)
			if parsedIP != "" {
				ip = parsedIP
				break
			}
		}
	}

	if strings.Contains(ip, ":") {
		host, _, err := net.SplitHostPort(ip)
		if err != nil {
			return ip
		}

		return host
	}

	return ip
}

func parseForwardedHeader(value string) string {
	parts := strings.Split(value, ",")
	parts = strings.Split(parts[0], ";")

	for _, part := range parts {
		kv := strings.Split(part, "=")

		if len(kv) == 2 {
			k := strings.ToLower(strings.TrimSpace(kv[0]))
			v := strings.TrimSpace(kv[1])

			if k == "for" {
				return v
			}
		}
	}

	return ""
}

func parseXForwardedForHeader(value string) string {
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[0])
}

func parseXRealIPHeader(value string) string {
	return value
}
