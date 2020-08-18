// Package tracking provides privacy-focused user tracking functions.
package tracking

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
)

// inspired by github.com/emvi/pirsch

// Tracker is the interface that wraps user tracking methods.
type Tracker interface {
	Delete(id string) error
	Get() ([]Hit, error)
	Hit(r *http.Request) error
	Searcher
}

// Searcher is the interface that wraps hits search methods.
type Searcher interface {
	Search(value string) ([]Hit, error)
	SearchByField(field, value string) ([]Hit, error)
}

// Hitter provides methods to hit requests and store them in a data store.
// In case of an error it will panic.
type Hitter struct {
	DB   *gorm.DB
	salt string
}

// NewTracker returns a new user tracker.
func NewTracker(db *gorm.DB, salt string) Tracker {
	return &Hitter{
		DB:   db,
		salt: salt,
	}
}

// Delete takes away the hit with the id specified from the database.
func (h *Hitter) Delete(id string) error {
	var hit Hit

	if err := h.DB.Delete(&hit, id).Error; err != nil {
		return fmt.Errorf("couldn't delete the hit: %v", err)
	}

	return nil
}

// Get lists all the hits stored in the database.
func (h *Hitter) Get() ([]Hit, error) {
	var hits []Hit

	if err := h.DB.Find(&hits).Error; err != nil {
		return nil, fmt.Errorf("couldn't find the hits: %v", err)
	}

	return hits, nil
}

// Hit stores the given request.
// The request might be ignored if it meets certain conditions.
func (h *Hitter) Hit(r *http.Request) error {
	if !ignoreHit(r) {
		hit, err := HitRequest(r, h.salt)
		if err != nil {
			return err
		}

		if err = h.DB.Create(&hit).Error; err != nil {
			return fmt.Errorf("couldn't save the hit: %v", err)
		}
	}

	return nil
}

// Search looks for a value and returns a slice of the hits that contain that value.
func (h *Hitter) Search(value string) ([]Hit, error) {
	var hits []Hit

	if err := h.DB.
		Where(`to_tsvector(id || ' ' || footprint || ' ' || 
		path || ' ' || url || ' ' || 
		language || ' ' || referer || ' ' || 
		user_agent || ' ' || date) @@ to_tsquery(?)`, value).
		Find(&hits).Error; err != nil {
		return nil, errors.New("no hits found")
	}

	return hits, nil
}

// SearchByField looks for a value and returns a slice of the hits that contain that value.
func (h *Hitter) SearchByField(field, value string) ([]Hit, error) {
	var hits []Hit

	if err := h.DB.Where("CONTAINS(?, '?')", field, value).Error; err != nil {
		return nil, errors.New("no hits found")
	}

	return hits, nil
}

// Check headers commonly used by bots.
// If the user is a bot return true, else return false.
func ignoreHit(r *http.Request) bool {
	if r.Header.Get("X-Moz") == "prefetch" ||
		r.Header.Get("X-Purpose") == "prefetch" ||
		r.Header.Get("X-Purpose") == "preview" ||
		r.Header.Get("Purpose") == "prefetch" ||
		r.Header.Get("Purpose") == "preview" {
		return true
	}

	userAgent := strings.ToLower(r.Header.Get("User-Agent"))

	for _, botUserAgent := range userAgentBotlist {
		if strings.Contains(userAgent, botUserAgent) {
			return true
		}

	}

	return false
}
