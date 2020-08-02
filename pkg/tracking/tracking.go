// Package tracking provides privacy-focused user tracking functions.
package tracking

import (
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// inspired on github.com/emvi/pirsch

// Hitter is the interface that wraps Hit methods.
type Hitter interface {
	Hit(r *http.Request) error
	DeleteHit(id string) error
	Get() ([]Hit, error)
}

// Searcher is the interface that wraps the basic search method.
type Searcher interface {
	Search(value string) ([]interface{}, error)
}

// Tracker provides methods to track requests and store them in a data store.
// In case of an error it will panic.
type Tracker struct {
	DB   *gorm.DB
	salt string
}

// NewTracker returns a new user tracker.
func NewTracker(db *gorm.DB, salt string) *Tracker {
	return &Tracker{
		DB:   db,
		salt: salt,
	}
}

// DeleteHit takes away the hit with the id specified from the database.
func (t *Tracker) DeleteHit(id string) error {
	var hit Hit
	err := t.DB.Delete(&hit, id).Error
	if err != nil {
		return errors.Wrap(err, "couldn't delete the hit")
	}

	return nil
}

// Get lists all the hits stored in the database.
func (t *Tracker) Get() ([]Hit, error) {
	var hits []Hit

	err := t.DB.Find(&hits).Error
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find the hits")
	}

	return hits, nil
}

// Hit stores the given request.
// The request might be ignored if it meets certain conditions.
func (t *Tracker) Hit(r *http.Request) error {
	if !ignoreHit(r) {
		hit := HitRequest(r, t.salt)

		err := t.DB.Create(&hit).Error
		if err != nil {
			return errors.Wrap(err, "couldn't save the hit")
		}
	}

	return nil
}

// Search looks for a value and returns a slice of the hits that contain that value.
func (t *Tracker) Search(value string) ([]Hit, error) {
	var hits []Hit

	err := t.DB.
		Where(`to_tsvector(id || ' ' || footprint || ' ' || 
		path || ' ' || url || ' ' || 
		language || ' ' || referer || ' ' || 
		user_agent || ' ' || date) @@ to_tsquery(?)`, value).
		Find(&hits).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "no hits found")
	}

	return hits, nil
}

// SearchByField looks for a value and returns a slice of the hits that contain that value.
func (t *Tracker) SearchByField(field, value string) ([]Hit, error) {
	var hits []Hit

	err := t.DB.Where("CONTAINS(?, '?')", field, value).Error
	if err != nil {
		return nil, errors.Wrap(err, "no hits found")
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
