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
		hits := t.FindAll()

		response.JSON(w, r, http.StatusOK, hits)
	}
}

// FindHitsByID prints the hit with the specified id.
func FindHitsByID(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		hit, _ := strconv.Atoi(id)

		hits := t.FindByID(int64(hit))

		response.JSON(w, r, http.StatusOK, hits)
	}
}

// FindHitsByDay prints the hit with the specified day.
func FindHitsByDay(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		day := mux.Vars(r)["day"]
		hit, _ := strconv.Atoi(day)

		hits := t.FindByDay(hit)

		response.HTMLText(w, r, http.StatusOK, hits)
	}
}

// FindHitsByHour prints the hit with the specified hour.
func FindHitsByHour(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hour := mux.Vars(r)["hour"]
		hit, _ := strconv.Atoi(hour)

		hits := t.FindByHour(hit)

		response.HTMLText(w, r, http.StatusOK, hits)
	}
}

// FindHitsByLanguage prints the hit with the specified language.
func FindHitsByLanguage(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		language := mux.Vars(r)["language"]

		var lang string

		switch language {
		case "english":
			lang = "en-US"
		case "spanish":
			lang = "es-ES"
		case "chinese":
			lang = "zh-CN"
		case "portuguese":
			lang = "pt-BR"
		case "german":
			lang = "de"
		case "french":
			lang = "fr"
		case "italian":
			lang = "it"
		case "russian":
			lang = "ru"
		case "hindi":
			lang = "in"
		case "japanese":
			lang = "jp"
		default:
			response.HTMLText(w, r, http.StatusBadRequest, "No hits found with this language")
			return
		}

		hits := t.FindByLanguage(lang)

		response.HTMLText(w, r, http.StatusOK, hits)
	}
}

// FindHitsByMonth prints the hit with the specified month.
func FindHitsByMonth(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		month := mux.Vars(r)["month"]
		hit, _ := strconv.Atoi(month)

		hits := t.FindByMonth(hit)

		response.HTMLText(w, r, http.StatusOK, hits)
	}
}

// FindHitsByPath prints the hit with the specified path.
func FindHitsByPath(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)["path"]

		hits := t.FindByPath(path)

		response.HTMLText(w, r, http.StatusOK, hits)
	}
}

// FindHitsByWeekday prints the hit with the specified weekday.
func FindHitsByWeekday(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		weekday := mux.Vars(r)["weekday"]
		hit, _ := strconv.Atoi(weekday)

		hits := t.FindByWeekday(hit)

		response.HTMLText(w, r, http.StatusOK, hits)
	}
}

// FindHitsByYear prints the hit with the specified year.
func FindHitsByYear(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year := mux.Vars(r)["year"]
		hit, _ := strconv.Atoi(year)

		hits := t.FindByYear(hit)

		response.HTMLText(w, r, http.StatusOK, hits)
	}
}
