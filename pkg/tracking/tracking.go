// Package tracking is a privacy-focused user tracker based on
// github.com/emvi/pirsch
package tracking

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Tracer is the interface that wraps the basic Hit method.
type Tracer interface {
	Hit(r *http.Request) error
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
func (tracker *Tracker) DeleteHit(id int64) int64 {
	var hit Hit
	total := tracker.DB.Where("id=?", id).Delete(&hit).RowsAffected

	return total
}

// FindAll lists all the hits stored in the database.
func (tracker *Tracker) FindAll() int64 {
	var hit Hit
	total := tracker.DB.Find(&hit).RowsAffected

	return total
}

// FindByID lists the hit with the id specified.
func (tracker *Tracker) FindByID(id int64) int64 {
	var hit Hit
	total := tracker.DB.Where("id=?", id).Find(&hit).RowsAffected

	return total
}

// FindByDay lists the hits that were stored in the day specified.
func (tracker *Tracker) FindByDay(day int) string {
	var hit Hit
	hits := tracker.DB.Where("day=?", day).Find(&hit).RowsAffected
	average := tracker.calculatePercentage(hits)

	result := fmt.Sprintf("<strong>Hits</strong>: %d<br><br><strong>Percentage of the total</strong>: %f%%", hits, average)

	return result
}

// FindByHour lists the hits that were stored in the hour specified.
func (tracker *Tracker) FindByHour(hour int) string {
	var hit Hit
	hits := tracker.DB.Where("hour=?", hour).Find(&hit).RowsAffected
	average := tracker.calculatePercentage(hits)

	result := fmt.Sprintf("<strong>Hits</strong>: %d<br><br><strong>Percentage of the total</strong>: %f%%", hits, average)

	return result
}

// FindByLanguage lists the hits that have the language specified.
func (tracker *Tracker) FindByLanguage(language string) string {
	var hit Hit
	hits := tracker.DB.Where("language=?", language).Find(&hit).RowsAffected
	average := tracker.calculatePercentage(hits)

	result := fmt.Sprintf("<strong>Hits</strong>: %d<br><br><strong>Percentage of the total</strong>: %f%%", hits, average)

	return result
}

// FindByMonth lists the hits that were stored in the month specified.
func (tracker *Tracker) FindByMonth(month int) string {
	var hit Hit
	hits := tracker.DB.Where("month=?", month).Find(&hit).RowsAffected
	average := tracker.calculatePercentage(hits)

	result := fmt.Sprintf("<strong>Hits</strong>: %d<br><br><strong>Percentage of the total</strong>: %f%%", hits, average)

	return result
}

// FindByPath lists the hits that were stored in the month specified.
func (tracker *Tracker) FindByPath(path string) string {
	var hit Hit
	hits := tracker.DB.Where("path=?", path).Find(&hit).RowsAffected
	average := tracker.calculatePercentage(hits)

	result := fmt.Sprintf("<strong>Hits</strong>: %d<br><br><strong>Percentage of the total</strong>: %f%%", hits, average)

	return result
}

// FindByWeekday lists the hits that were stored in the weekday specified.
func (tracker *Tracker) FindByWeekday(weekday int) string {
	var hit Hit
	hits := tracker.DB.Where("weekday=?", weekday).Find(&hit).RowsAffected
	average := tracker.calculatePercentage(hits)

	result := fmt.Sprintf("<strong>Hits</strong>: %d<br><br><strong>Percentage of the total</strong>: %f%%", hits, average)

	return result
}

// FindByYear lists the hits that were stored in the year specified.
func (tracker *Tracker) FindByYear(year int) string {
	var hit Hit
	hits := tracker.DB.Where("year=?", year).Find(&hit).RowsAffected
	average := tracker.calculatePercentage(hits)

	result := fmt.Sprintf("<strong>Hits</strong>: %d<br><br><strong>Percentage of the total</strong>: %f%%", hits, average)

	return result
}

// Hit stores the given request.
// The request might be ignored if it meets certain conditions.
// The actions performed within this function run in their own goroutine, so you don't need to create one yourself.
func (tracker *Tracker) Hit(r *http.Request) error {
	if !ignoreHit(r) {
		hit := HitRequest(r, tracker.salt)

		err := tracker.DB.Create(hit).Error
		if err != nil {
			return errors.Wrap(err, "couldn't create the hit")
		}
	}

	return nil
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

// calculateAverage returns the percentage that hits represent of the total hits stored
func (tracker *Tracker) calculatePercentage(hits int64) float64 {
	var hit Hit
	total := tracker.DB.Find(&hit).RowsAffected

	percentage := (100 / (float64(total)) * float64(hits))

	return percentage
}
