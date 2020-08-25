package tracking

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/GGP1/palo/internal/random"
)

// Hit represents a single data point/page visit.
type Hit struct {
	ID        string    `json:"id"`
	Footprint string    `json:"footprint"`
	Path      string    `json:"path"`
	URL       string    `json:"url"`
	Language  string    `json:"language"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Referer   string    `json:"referer"`
	Date      time.Time `json:"date"`
}

// String returns a string of the hit struct.
func (hit *Hit) String() (string, error) {
	out, err := json.Marshal(hit)
	if err != nil {
		return "", errors.New("couldn't marshal the hit")
	}

	return string(out), nil
}

// HitRequest generates a hit for each request.
func HitRequest(r *http.Request, salt string) (*Hit, error) {
	id := random.GenerateRunes(27)
	date := time.Now()

	footprint, err := Footprint(r, salt)
	if err != nil {
		return nil, err
	}

	return &Hit{
		ID:        id,
		Footprint: footprint,
		Path:      r.URL.Path,
		URL:       r.URL.String(),
		Language:  getLanguage(r),
		UserAgent: r.UserAgent(),
		Referer:   r.Header.Get("Referer"),
		Date:      date,
	}, nil
}

// Get the user language from the web browser.
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
			language = "english"
		case spanish:
			language = "spanish"
		case portuguese:
			language = "portuguese"
		case chinese:
			language = "chinese"
		case german:
			language = "german"
		case french:
			language = "french"
		case italian:
			language = "italian"
		case russian:
			language = "russian"
		case hindi:
			language = "hindi"
		case japanese:
			language = "japanese"
		}
		return language
	}

	return ""
}
