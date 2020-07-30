package tracking

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Hit represents a single data point/page visit.
type Hit struct {
	ID        string    `json:"id"`
	Footprint string    `json:"footprint"`
	Path      string    `json:"path"`
	URL       string    `json:"url"`
	Language  string    `json:"language"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	Date      time.Time `json:"date"`
}

// String returns a string of the hit struct
func (hit *Hit) String() string {
	out, _ := json.Marshal(hit)
	return string(out)
}

// HitRequest generates a hit for each request
func HitRequest(r *http.Request, salt string) *Hit {
	id := uuid.New()
	date := time.Now()

	return &Hit{
		ID:        id.String(),
		Footprint: Footprint(r, salt),
		Path:      r.URL.Path,
		URL:       r.URL.String(),
		Language:  getLanguage(r),
		UserAgent: r.UserAgent(),
		Referer:   r.Header.Get("Referer"),
		Date:      date,
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
