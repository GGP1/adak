package tracking

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Hit represents a single data point/page visit.
type Hit struct {
	ID        int64        `json:"id"`
	Footprint string       `json:"footprint"`
	Path      string       `json:"path"`
	URL       string       `json:"url"`
	Language  string       `json:"language"`
	UserAgent string       `json:"user_agent"`
	Referer   string       `json:"referer"`
	Hour      int          `json:"hour"`
	Weekday   time.Weekday `json:"weekday"`
	Day       int          `json:"day"`
	Month     time.Month   `json:"month"`
	Year      int          `json:"year"`
}

// String returns a string of the hit struct
func (hit *Hit) String() string {
	out, _ := json.Marshal(hit)
	return string(out)
}

// HitRequest generates a hit for each request
func HitRequest(r *http.Request, salt string) *Hit {
	hour := time.Now().Hour()
	y, m, d := time.Now().Date()
	weekday := time.Now().Weekday()

	return &Hit{
		Footprint: Footprint(r, salt),
		Path:      r.URL.Path,
		URL:       r.URL.String(),
		Language:  getLanguage(r),
		UserAgent: r.UserAgent(),
		Referer:   r.Header.Get("Referer"),
		Hour:      hour,
		Weekday:   weekday,
		Day:       d,
		Month:     m,
		Year:      y,
	}
}

// Get the user language
func getLanguage(r *http.Request) string {
	lang := r.Header.Get("Accept-Language")

	if lang != "" {
		langs := strings.Split(lang, ";")
		parts := strings.Split(langs[0], ",")
		return parts[0]
	}

	return ""
}
