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
		english := "en-US"
		spanish := "es-ES"
		chinese := "zh-CN"
		portuguese := "pt-BR"
		german := "de"
		french := "fr"
		italian := "it"
		russian := "ru"
		hindi := "in"
		japanese := "jp"

		langs := strings.Split(lang, ";")
		parts := strings.Split(langs[0], ",")
		var language string

		switch parts[0] {
		case english:
			language = "English"
		case spanish:
			language = "Spanish"
		case portuguese:
			language = "Portuguese"
		case chinese:
			language = "Chinese"
		case german:
			language = "German"
		case french:
			language = "French"
		case italian:
			language = "Italian"
		case russian:
			language = "Russian"
		case hindi:
			language = "Hindi"
		case japanese:
			language = "Japanese"
		}
		return language
	}

	return ""
}