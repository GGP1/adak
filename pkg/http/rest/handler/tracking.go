package handler

import (
	"net/http"
	"strconv"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/tracking"
	"github.com/gorilla/mux"
)

// FindHits retrieves total amount of hits stored.
func FindHits(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		total := t.FindAll()

		response.JSON(w, r, http.StatusOK, total)
	}
}

// FindHitsByID prints the hit with the specified id.
func FindHitsByID(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		hit, _ := strconv.Atoi(id)

		total := t.FindByID(int64(hit))

		response.JSON(w, r, http.StatusOK, total)
	}
}

// FindHitsByDay prints the hit with the specified day.
func FindHitsByDay(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		day := mux.Vars(r)["day"]
		hit, _ := strconv.Atoi(day)

		total := t.FindByDay(hit)

		response.JSON(w, r, http.StatusOK, total)
	}
}

// FindHitsByHour prints the hit with the specified hour.
func FindHitsByHour(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hour := mux.Vars(r)["hour"]
		hit, _ := strconv.Atoi(hour)

		total := t.FindByHour(hit)

		response.JSON(w, r, http.StatusOK, total)
	}
}

// FindHitsByLanguage prints the hit with the specified language.
func FindHitsByLanguage(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		language := mux.Vars(r)["language"]

		total := t.FindByLanguage(language)

		response.JSON(w, r, http.StatusOK, total)
	}
}

// FindHitsByMonth prints the hit with the specified month.
func FindHitsByMonth(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		month := mux.Vars(r)["month"]
		hit, _ := strconv.Atoi(month)

		total := t.FindByMonth(hit)

		response.JSON(w, r, http.StatusOK, total)
	}
}

// FindHitsByPath prints the hit with the specified path.
func FindHitsByPath(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)["path"]

		total := t.FindByPath(path)

		response.JSON(w, r, http.StatusOK, total)
	}
}

// FindHitsByWeekday prints the hit with the specified weekday.
func FindHitsByWeekday(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		weekday := mux.Vars(r)["weekday"]
		hit, _ := strconv.Atoi(weekday)

		total := t.FindByWeekday(hit)

		response.JSON(w, r, http.StatusOK, total)
	}
}

// FindHitsByYear prints the hit with the specified year.
func FindHitsByYear(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year := mux.Vars(r)["year"]
		hit, _ := strconv.Atoi(year)

		total := t.FindByYear(hit)

		response.JSON(w, r, http.StatusOK, total)
	}
}
